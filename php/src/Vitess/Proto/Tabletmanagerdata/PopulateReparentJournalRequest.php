<?php
// DO NOT EDIT! Generated by Protobuf-PHP protoc plugin 1.0
// Source: tabletmanagerdata.proto
//   Date: 2016-01-22 01:34:35

namespace Vitess\Proto\Tabletmanagerdata {

  class PopulateReparentJournalRequest extends \DrSlump\Protobuf\Message {

    /**  @var int */
    public $time_created_ns = null;
    
    /**  @var string */
    public $action_name = null;
    
    /**  @var \Vitess\Proto\Topodata\TabletAlias */
    public $master_alias = null;
    
    /**  @var string */
    public $replication_position = null;
    

    /** @var \Closure[] */
    protected static $__extensions = array();

    public static function descriptor()
    {
      $descriptor = new \DrSlump\Protobuf\Descriptor(__CLASS__, 'tabletmanagerdata.PopulateReparentJournalRequest');

      // OPTIONAL INT64 time_created_ns = 1
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 1;
      $f->name      = "time_created_ns";
      $f->type      = \DrSlump\Protobuf::TYPE_INT64;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      // OPTIONAL STRING action_name = 2
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 2;
      $f->name      = "action_name";
      $f->type      = \DrSlump\Protobuf::TYPE_STRING;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      // OPTIONAL MESSAGE master_alias = 3
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 3;
      $f->name      = "master_alias";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $f->reference = '\Vitess\Proto\Topodata\TabletAlias';
      $descriptor->addField($f);

      // OPTIONAL STRING replication_position = 4
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 4;
      $f->name      = "replication_position";
      $f->type      = \DrSlump\Protobuf::TYPE_STRING;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      foreach (self::$__extensions as $cb) {
        $descriptor->addField($cb(), true);
      }

      return $descriptor;
    }

    /**
     * Check if <time_created_ns> has a value
     *
     * @return boolean
     */
    public function hasTimeCreatedNs(){
      return $this->_has(1);
    }
    
    /**
     * Clear <time_created_ns> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function clearTimeCreatedNs(){
      return $this->_clear(1);
    }
    
    /**
     * Get <time_created_ns> value
     *
     * @return int
     */
    public function getTimeCreatedNs(){
      return $this->_get(1);
    }
    
    /**
     * Set <time_created_ns> value
     *
     * @param int $value
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function setTimeCreatedNs( $value){
      return $this->_set(1, $value);
    }
    
    /**
     * Check if <action_name> has a value
     *
     * @return boolean
     */
    public function hasActionName(){
      return $this->_has(2);
    }
    
    /**
     * Clear <action_name> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function clearActionName(){
      return $this->_clear(2);
    }
    
    /**
     * Get <action_name> value
     *
     * @return string
     */
    public function getActionName(){
      return $this->_get(2);
    }
    
    /**
     * Set <action_name> value
     *
     * @param string $value
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function setActionName( $value){
      return $this->_set(2, $value);
    }
    
    /**
     * Check if <master_alias> has a value
     *
     * @return boolean
     */
    public function hasMasterAlias(){
      return $this->_has(3);
    }
    
    /**
     * Clear <master_alias> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function clearMasterAlias(){
      return $this->_clear(3);
    }
    
    /**
     * Get <master_alias> value
     *
     * @return \Vitess\Proto\Topodata\TabletAlias
     */
    public function getMasterAlias(){
      return $this->_get(3);
    }
    
    /**
     * Set <master_alias> value
     *
     * @param \Vitess\Proto\Topodata\TabletAlias $value
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function setMasterAlias(\Vitess\Proto\Topodata\TabletAlias $value){
      return $this->_set(3, $value);
    }
    
    /**
     * Check if <replication_position> has a value
     *
     * @return boolean
     */
    public function hasReplicationPosition(){
      return $this->_has(4);
    }
    
    /**
     * Clear <replication_position> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function clearReplicationPosition(){
      return $this->_clear(4);
    }
    
    /**
     * Get <replication_position> value
     *
     * @return string
     */
    public function getReplicationPosition(){
      return $this->_get(4);
    }
    
    /**
     * Set <replication_position> value
     *
     * @param string $value
     * @return \Vitess\Proto\Tabletmanagerdata\PopulateReparentJournalRequest
     */
    public function setReplicationPosition( $value){
      return $this->_set(4, $value);
    }
  }
}

