// Copyright 2015, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vtgate

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/net/context"

	mproto "github.com/youtube/vitess/go/mysql/proto"
	"github.com/youtube/vitess/go/stats"
	"github.com/youtube/vitess/go/vt/discovery"
	pbq "github.com/youtube/vitess/go/vt/proto/query"
	pbt "github.com/youtube/vitess/go/vt/proto/topodata"
	"github.com/youtube/vitess/go/vt/proto/vtrpc"
	tproto "github.com/youtube/vitess/go/vt/tabletserver/proto"
	"github.com/youtube/vitess/go/vt/tabletserver/tabletconn"
	"github.com/youtube/vitess/go/vt/topo"
	"github.com/youtube/vitess/go/vt/vterrors"
)

var (
	cellsToWatch        = flag.String("cells_to_watch", "", "comma-separated list of cells for watching endpoints")
	refreshInterval     = flag.Duration("endpoint_refresh_interval", 1*time.Minute, "endpoint refresh interval")
	topoReadConcurrency = flag.Int("topo_read_concurrency", 32, "concurrent topo reads")
)

const gatewayImplementationDiscovery = "discoverygateway"

func init() {
	RegisterGatewayCreator(gatewayImplementationDiscovery, createDiscoveryGateway)
}

func createDiscoveryGateway(hc discovery.HealthCheck, topoServer topo.Server, serv SrvTopoServer, cell string, retryDelay time.Duration, retryCount int, connTimeoutTotal, connTimeoutPerConn, connLife time.Duration, connTimings *stats.MultiTimings) Gateway {
	return &discoveryGateway{
		hc:              hc,
		topoServer:      topoServer,
		localCell:       cell,
		retryCount:      retryCount,
		tabletsWatchers: make([]*discovery.CellTabletsWatcher, 0, 1),
	}
}

type discoveryGateway struct {
	hc         discovery.HealthCheck
	topoServer topo.Server
	localCell  string
	retryCount int

	tabletsWatchers []*discovery.CellTabletsWatcher
}

// InitializeConnections creates connections to VTTablets.
func (dg *discoveryGateway) InitializeConnections(ctx context.Context) error {
	dg.hc.SetListener(dg)
	for _, cell := range strings.Split(*cellsToWatch, ",") {
		ctw := discovery.NewCellTabletsWatcher(dg.topoServer, dg.hc, cell, *refreshInterval, *topoReadConcurrency)
		dg.tabletsWatchers = append(dg.tabletsWatchers, ctw)
	}
	return nil
}

// Execute executes the non-streaming query for the specified keyspace, shard, and tablet type.
func (dg *discoveryGateway) Execute(ctx context.Context, keyspace, shard string, tabletType pbt.TabletType, query string, bindVars map[string]interface{}, transactionID int64) (qr *mproto.QueryResult, err error) {
	err = dg.withRetry(ctx, keyspace, shard, tabletType, func(conn tabletconn.TabletConn) error {
		var innerErr error
		qr, innerErr = conn.Execute2(ctx, query, bindVars, transactionID)
		return innerErr
	}, transactionID, false)
	return qr, err
}

// ExecuteBatch executes a group of queries for the specified keyspace, shard, and tablet type.
func (dg *discoveryGateway) ExecuteBatch(ctx context.Context, keyspace, shard string, tabletType pbt.TabletType, queries []tproto.BoundQuery, asTransaction bool, transactionID int64) (qrs *tproto.QueryResultList, err error) {
	err = dg.withRetry(ctx, keyspace, shard, tabletType, func(conn tabletconn.TabletConn) error {
		var innerErr error
		qrs, innerErr = conn.ExecuteBatch2(ctx, queries, asTransaction, transactionID)
		return innerErr
	}, transactionID, false)
	return qrs, err
}

// StreamExecute executes a streaming query for the specified keyspace, shard, and tablet type.
func (dg *discoveryGateway) StreamExecute(ctx context.Context, keyspace, shard string, tabletType pbt.TabletType, query string, bindVars map[string]interface{}, transactionID int64) (<-chan *mproto.QueryResult, tabletconn.ErrFunc) {
	var usedConn tabletconn.TabletConn
	var erFunc tabletconn.ErrFunc
	var results <-chan *mproto.QueryResult
	err := dg.withRetry(ctx, keyspace, shard, tabletType, func(conn tabletconn.TabletConn) error {
		var err error
		results, erFunc, err = conn.StreamExecute2(ctx, query, bindVars, transactionID)
		usedConn = conn
		return err
	}, transactionID, true)
	if err != nil {
		return results, func() error { return err }
	}
	inTransaction := (transactionID != 0)
	return results, func() error {
		return WrapError(erFunc(), keyspace, shard, tabletType, usedConn.EndPoint(), inTransaction)
	}
}

// Begin starts a transaction for the specified keyspace, shard, and tablet type.
// It returns the transaction ID.
func (dg *discoveryGateway) Begin(ctx context.Context, keyspace string, shard string, tabletType pbt.TabletType) (transactionID int64, err error) {
	err = dg.withRetry(ctx, keyspace, shard, tabletType, func(conn tabletconn.TabletConn) error {
		var innerErr error
		transactionID, innerErr = conn.Begin2(ctx)
		return innerErr
	}, 0, false)
	return transactionID, err
}

// Commit commits the current transaction for the specified keyspace, shard, and tablet type.
func (dg *discoveryGateway) Commit(ctx context.Context, keyspace, shard string, tabletType pbt.TabletType, transactionID int64) error {
	return dg.withRetry(ctx, keyspace, shard, tabletType, func(conn tabletconn.TabletConn) error {
		return conn.Commit2(ctx, transactionID)
	}, transactionID, false)
}

// Rollback rolls back the current transaction for the specified keyspace, shard, and tablet type.
func (dg *discoveryGateway) Rollback(ctx context.Context, keyspace, shard string, tabletType pbt.TabletType, transactionID int64) error {
	return dg.withRetry(ctx, keyspace, shard, tabletType, func(conn tabletconn.TabletConn) error {
		return conn.Rollback2(ctx, transactionID)
	}, transactionID, false)
}

// SplitQuery splits a query into sub-queries for the specified keyspace, shard, and tablet type.
func (dg *discoveryGateway) SplitQuery(ctx context.Context, keyspace, shard string, tabletType pbt.TabletType, sql string, bindVariables map[string]interface{}, splitColumn string, splitCount int) (queries []tproto.QuerySplit, err error) {
	err = dg.withRetry(ctx, keyspace, shard, tabletType, func(conn tabletconn.TabletConn) error {
		var innerErr error
		queries, innerErr = conn.SplitQuery(ctx, tproto.BoundQuery{
			Sql:           sql,
			BindVariables: bindVariables,
		}, splitColumn, splitCount)
		return innerErr
	}, 0, false)
	return
}

// Close shuts down underlying connections.
func (dg *discoveryGateway) Close(ctx context.Context) error {
	for _, ctw := range dg.tabletsWatchers {
		ctw.Stop()
	}
	return nil
}

// StatsUpdate receives updates about target and realtime stats changes.
func (dg *discoveryGateway) StatsUpdate(endPoint *pbt.EndPoint, cell, name string, target *pbq.Target, serving bool, tabletExternallyReparentedTimestamp int64, stats *pbq.RealtimeStats) {
}

// withRetry gets available connections and executes the action. If there are retryable errors,
// it retries retryCount times before failing. It does not retry if the connection is in
// the middle of a transaction. While returning the error check if it maybe a result of
// a resharding event, and set the re-resolve bit and let the upper layers
// re-resolve and retry.
func (dg *discoveryGateway) withRetry(ctx context.Context, keyspace, shard string, tabletType pbt.TabletType, action func(conn tabletconn.TabletConn) error, transactionID int64, isStreaming bool) error {
	var endPointLastUsed *pbt.EndPoint
	var err error
	inTransaction := (transactionID != 0)
	invalidEndPoints := make(map[string]bool)

	for i := 0; i < dg.retryCount+1; i++ {
		var endPoint *pbt.EndPoint
		endPoints := dg.getEndPoints(keyspace, shard, tabletType)
		if len(endPoints) == 0 {
			// fail fast if there is no endpoint
			err = vterrors.FromError(vtrpc.ErrorCode_INTERNAL_ERROR, fmt.Errorf("no valid endpoint"))
			break
		}
		shuffleEndPoints(endPoints)

		// skip endpoints we tried before
		for _, ep := range endPoints {
			if _, ok := invalidEndPoints[discovery.EndPointToMapKey(ep)]; !ok {
				endPoint = ep
				break
			}
		}
		if endPoint == nil {
			if err == nil {
				// do not override error from last attempt.
				err = vterrors.FromError(vtrpc.ErrorCode_INTERNAL_ERROR, fmt.Errorf("no available connection"))
			}
			break
		}

		// execute
		endPointLastUsed = endPoint
		conn := dg.hc.GetConnection(endPoint)
		if conn == nil {
			err = vterrors.FromError(vtrpc.ErrorCode_INTERNAL_ERROR, fmt.Errorf("no connection for %+v", endPoint))
			invalidEndPoints[discovery.EndPointToMapKey(endPoint)] = true
			continue
		}
		err = action(conn)
		if dg.canRetry(ctx, err, transactionID, isStreaming) {
			invalidEndPoints[discovery.EndPointToMapKey(endPoint)] = true
			continue
		}
		break
	}
	return WrapError(err, keyspace, shard, tabletType, endPointLastUsed, inTransaction)
}

// canRetry determines whether a query can be retried or not.
// OperationalErrors like retry/fatal are retryable if query is not in a txn.
// All other errors are non-retryable.
func (dg *discoveryGateway) canRetry(ctx context.Context, err error, transactionID int64, isStreaming bool) bool {
	if err == nil {
		return false
	}
	// Do not retry if ctx.Done() is closed.
	select {
	case <-ctx.Done():
		return false
	default:
	}
	if serverError, ok := err.(*tabletconn.ServerError); ok {
		switch serverError.Code {
		case tabletconn.ERR_FATAL:
			// Do not retry on fatal error for streaming query.
			// For streaming query, vttablet sends:
			// - RETRY, if streaming is not started yet;
			// - FATAL, if streaming is broken halfway.
			// For non-streaming query, handle as ERR_RETRY.
			if isStreaming {
				return false
			}
			fallthrough
		case tabletconn.ERR_RETRY:
			// Retry on RETRY and FATAL if not in a transaction.
			inTransaction := (transactionID != 0)
			return !inTransaction
		default:
			// Not retry for TX_POOL_FULL and normal server errors.
			return false
		}
	}
	// Do not retry on operational error.
	return false
}

func shuffleEndPoints(endPoints []*pbt.EndPoint) {
	index := 0
	length := len(endPoints)
	for i := length - 1; i > 0; i-- {
		index = rand.Intn(i + 1)
		endPoints[i], endPoints[index] = endPoints[index], endPoints[i]
	}
}

// getEndPoints gets all available endpoints from HealthCheck,
// and selects the usable ones based several rules:
// master - return one from any cells with latest reparent timestamp;
// replica - return all from local cell.
// TODO(liang): select replica by replication lag.
func (dg *discoveryGateway) getEndPoints(keyspace, shard string, tabletType pbt.TabletType) []*pbt.EndPoint {
	epsList := dg.hc.GetEndPointStatsFromTarget(keyspace, shard, tabletType)
	// for master, use any cells and return the one with max reparent timestamp.
	if tabletType == pbt.TabletType_MASTER {
		var maxTimestamp int64
		var ep *pbt.EndPoint
		for _, eps := range epsList {
			if eps.LastError != nil || !eps.Serving {
				continue
			}
			if eps.TabletExternallyReparentedTimestamp >= maxTimestamp {
				maxTimestamp = eps.TabletExternallyReparentedTimestamp
				ep = eps.EndPoint
			}
		}
		if ep == nil {
			return nil
		}
		return []*pbt.EndPoint{ep}
	}
	// for non-master, use only endpoints from local cell.
	var epList []*pbt.EndPoint
	for _, eps := range epsList {
		if eps.LastError != nil || !eps.Serving {
			continue
		}
		if dg.localCell != eps.Cell {
			continue
		}
		epList = append(epList, eps.EndPoint)
	}
	return epList
}

// WrapError returns ShardConnError which preserves the original error code if possible,
// adds the connection context
// and adds a bit to determine whether the keyspace/shard needs to be
// re-resolved for a potential sharding event.
func WrapError(in error, keyspace, shard string, tabletType pbt.TabletType, endPoint *pbt.EndPoint, inTransaction bool) (wrapped error) {
	if in == nil {
		return nil
	}
	shardIdentifier := fmt.Sprintf("%s.%s.%s, %+v", keyspace, shard, strings.ToLower(tabletType.String()), endPoint)
	code := tabletconn.ERR_NORMAL
	serverError, ok := in.(*tabletconn.ServerError)
	if ok {
		code = serverError.Code
	}

	shardConnErr := &ShardConnError{
		Code:            code,
		ShardIdentifier: shardIdentifier,
		InTransaction:   inTransaction,
		Err:             in,
		EndPointCode:    vterrors.RecoverVtErrorCode(in),
	}
	return shardConnErr
}
