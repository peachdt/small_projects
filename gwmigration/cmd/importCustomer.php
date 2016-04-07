<?php

if(php_sapi_name()!=="cli"){
    echo "Must be run from the commend line.";
};

/**
 * Setup a magento instance so we can run this export from the command line.
 */

require_once('/var/www/magento/app/Mage.php');
umask(0);

if (!Mage::isInstalled()) {
    echo "Application is not installed yet, please complete install wizard first.";
    exit;
}

main();

function readFromCSV($path)
{
    $file = fopen($path, 'r');
    while (($line = fgetcsv($file)) !== FALSE)
    {
        if (!empty($line[0]) && !empty($line[1]))
        {
            $data['first_name'] = $line[0];
            $data['last_name'] = $line[1];
            $data['email'] = $line[2];
            // use random password
            $data['password'] = generatePassword();
            // take out the timezone since all dates are in UTC
            // 2015-06-08 21:22:53 +0000
            $temp_date = explode(" ", $line[3]);
            $stripped_date = $temp_date[0] . " " . $temp_date[1];
            $data['created_date'] = $stripped_date;
            upsertCustomer($data);
            unset($data);
        }

    }
}

function getRandomBytes($nbBytes = 32)
{
    $bytes = openssl_random_pseudo_bytes($nbBytes, $strong);
    if (false !== $bytes && true === $strong) {
        return $bytes;
    }
    else {
        throw new \Exception("Unable to generate secure token from OpenSSL.");
    }
}

function generatePassword(){
   // random length from 80 to 1000
    $length = rand(80,100);
     return substr(preg_replace("/[^a-zA-Z0-9]/", "", base64_encode(getRandomBytes($length+1))),0,$length);
}

function upsertCustomer($data)
{
    $websiteId = Mage::app()->getWebsite()->getId();
    $store = Mage::app()->getStore();
    $customer = Mage::getModel("customer/customer");
    $customer->setWebsiteId($websiteId);
    $customer->setStore($store);
    $customer->loadByEmail($data['email']);
    $customer->setFirstname($data['first_name'])
        ->setLastname($data['last_name'])
        ->setEmail($data['email'])
        ->setPassword($data['password'])
        ->setCreatedAt($data['created_date']);
    try {
        $customer->save();
    } catch (Exception $e) {
        Zend_Debug::dump($e->getMessage());
    }

}

function main() {
    // path to csv file
    $path = "./var/import/BC_customers.csv";

    $time_start = microtime(true);
    echo "timer starts" . "\n";
    readFromCSV($path);
    $time_end = microtime(true);

//dividing with 60 will give the execution time in minutes other wise seconds
    $execution_time = ($time_end - $time_start)/60;

//execution time of the script
    echo 'Total Execution Time:'.$execution_time.' Mins' . "\n";
}