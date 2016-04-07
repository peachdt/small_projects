<?php

class Format
{
    public static function toUpper($text)
    {
        return strtoupper($text);
    }

    public static function toLower($text)
    {
        return strtolower($text);
    }

    function toProper($text)
    {
      return self::toUpper(substr($text, 0, 1)) . self::toLower(substr($text, 1));
    }
}
