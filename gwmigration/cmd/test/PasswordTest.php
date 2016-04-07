<?php ini_set("memory_limit","256M");

require_once dirname(__FILE__).'/../../cmd/test_pw_generator.php';

class PasswordTest extends PHPUnit_Framework_TestCase {

    public function testToProper() {
        $pw_list = array();
        for ($i = 0; $i < 1000000; $i++) {
            $password = GeneratePW::generatePassword();
            $this->assertTrue(strlen(GeneratePW::generatePassword()) > 60);
            $this->assertFalse(array_key_exists($password, $pw_list));
            $pw_list[$password] = "";
        }
    }
}
