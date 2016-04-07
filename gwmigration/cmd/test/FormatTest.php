<?php
//require_once 'PHPUnit/Framework.php';
require_once dirname(__FILE__).'/../../cmd/format.php';

class FormatTest extends PHPUnit_Framework_TestCase {

  public function testToProper() {
    $this->assertEquals('SEbastian', Format::toProper( 'sebastian' ));
  }
}
