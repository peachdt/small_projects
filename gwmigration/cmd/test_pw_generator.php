<?php

class GeneratePW
{
// generate 10 million passwords and check if any of them are duplicates
    public static function getRandomBytes($nbBytes = 32)
    {
        $bytes = openssl_random_pseudo_bytes($nbBytes, $strong);
        if (false !== $bytes && true === $strong) {
            return $bytes;
        } else {
            throw new \Exception("Unable to generate secure token from OpenSSL.");
        }
    }

    public static function generatePassword()
    {
        $length = rand(80, 100);
        return substr(preg_replace("/[^a-zA-Z0-9]/", "", base64_encode(self::getRandomBytes($length + 1))), 0, $length);
    }
}