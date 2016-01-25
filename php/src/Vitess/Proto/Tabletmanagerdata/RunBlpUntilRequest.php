<?php
// DO NOT EDIT! Generated by Protobuf-PHP protoc plugin 1.0
// Source: tabletmanagerdata.proto
//   Date: 2016-01-22 01:34:35

namespace Vitess\Proto\Tabletmanagerdata {

  class RunBlpUntilRequest extends \DrSlump\Protobuf\Message {

    /**  @var \Vitess\Proto\Tabletmanagerdata\BlpPosition[]  */
    public $blp_positions = array();
    
    /**  @var int */
    public $wait_timeout = null;
    

    /** @var \Closure[] */
    protected static $__extensions = array();

    public static function descriptor()
    {
      $descriptor = new \DrSlump\Protobuf\Descriptor(__CLASS__, 'tabletmanagerdata.RunBlpUntilRequest');

      // REPEATED MESSAGE blp_positions = 1
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 1;
      $f->name      = "blp_positions";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_REPEATED;
      $f->reference = '\Vitess\Proto\Tabletmanagerdata\BlpPosition';
      $descriptor->addField($f);

      // OPTIONAL INT64 wait_timeout = 2
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 2;
      $f->name      = "wait_timeout";
      $f->type      = \DrSlump\Protobuf::TYPE_INT64;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      foreach (self::$__extensions as $cb) {
        $descriptor->addField($cb(), true);
      }

      return $descriptor;
    }

    /**
     * Check if <blp_positions> has a value
     *
     * @return boolean
     */
    public function hasBlpPositions(){
      return $this->_has(1);
    }
    
    /**
     * Clear <blp_positions> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\RunBlpUntilRequest
     */
    public function clearBlpPositions(){
      return $this->_clear(1);
    }
    
    /**
     * Get <blp_positions> value
     *
     * @param int $idx
     * @return \Vitess\Proto\Tabletmanagerdata\BlpPosition
     */
    public function getBlpPositions($idx = NULL){
      return $this->_get(1, $idx);
    }
    
    /**
     * Set <blp_positions> value
     *
     * @param \Vitess\Proto\Tabletmanagerdata\BlpPosition $value
     * @return \Vitess\Proto\Tabletmanagerdata\RunBlpUntilRequest
     */
    public function setBlpPositions(\Vitess\Proto\Tabletmanagerdata\BlpPosition $value, $idx = NULL){
      return $this->_set(1, $value, $idx);
    }
    
    /**
     * Get all elements of <blp_positions>
     *
     * @return \Vitess\Proto\Tabletmanagerdata\BlpPosition[]
     */
    public function getBlpPositionsList(){
     return $this->_get(1);
    }
    
    /**
     * Add a new element to <blp_positions>
     *
     * @param \Vitess\Proto\Tabletmanagerdata\BlpPosition $value
     * @return \Vitess\Proto\Tabletmanagerdata\RunBlpUntilRequest
     */
    public function addBlpPositions(\Vitess\Proto\Tabletmanagerdata\BlpPosition $value){
     return $this->_add(1, $value);
    }
    
    /**
     * Check if <wait_timeout> has a value
     *
     * @return boolean
     */
    public function hasWaitTimeout(){
      return $this->_has(2);
    }
    
    /**
     * Clear <wait_timeout> value
     *
     * @return \Vitess\Proto\Tabletmanagerdata\RunBlpUntilRequest
     */
    public function clearWaitTimeout(){
      return $this->_clear(2);
    }
    
    /**
     * Get <wait_timeout> value
     *
     * @return int
     */
    public function getWaitTimeout(){
      return $this->_get(2);
    }
    
    /**
     * Set <wait_timeout> value
     *
     * @param int $value
     * @return \Vitess\Proto\Tabletmanagerdata\RunBlpUntilRequest
     */
    public function setWaitTimeout( $value){
      return $this->_set(2, $value);
    }
  }
}

