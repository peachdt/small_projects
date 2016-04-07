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
Mage::init();
main();

function readFromCSV($path)
{
    $file = fopen($path, 'r');
    while (($line = fgetcsv($file)) !== FALSE)
    {
        if (!empty($line[0]) && !empty($line[1]))
        {
            $data['first_name'] = $line[1];
            $data['last_name'] = $line[2];
            $data['company'] = $line[3];
            $data['street_1'] = $line[4];
            $data['street_2'] = $line[5];
            $data['city'] = $line[6];
            $data['state'] = $line[7];
            $data['state_iso2'] = $line[8];
            $data['zip_code'] = $line[9];
            $data['country'] = $line[10];
            $data['phone'] = $line[11];
            $data['email'] = $line[12];
            upsertAddress($data);
            unset($data);
        }

    }
}

function upsertAddress($data)
{
    $websiteId = Mage::app()->getWebsite()->getId();
    $store = Mage::app()->getStore();
    $customer = Mage::getModel("customer/customer");
    $customer->setWebsiteId($websiteId);
    $customer->setStore($store);
    $customer->loadByEmail($data['email']);

    if ($customer->getId()){
        $address = Mage::getModel("customer/address");
        $address->setCustomerId($customer->getId())
            ->setFirstname($data['first_name'])
            ->setLastname($data['last_name'])
            ->setCompant($data['company'])
            ->setTelephone($data['phone'])
            ->setStreetFull(array($data['street_1'], $data['street_2']))
            ->setCity($data['city'])
            ->setCountryId($data['country'])
            ->setPostcode($data['zip_code']);
            $region = Mage::getModel('directory/region')->loadByCode($data['state_iso2'], $data['country']);
            $state_id = $region->getId();
            if ($region->getId()) {
                $address->setRegionId($state_id);
            } else {
                if ($data['state'] !== '') {
                    $address->setRegion($data['state']);
                }
            }
            $address->setIsDefaultBilling('1')
            ->setIsDefaultShipping('1')
            ->setSaveInAddressBook('1');

        try {
            $address->save();
        } catch (Exception $e) {
            Zend_Debug::dump($e->getMessage());
        }
    }
}

function main() {
    // path to csv file
    $path = "./var/import/BC_address.csv";

    $time_start = microtime(true);
    echo "timer starts" . "\n";
    echo "reading file from " . $path;
    readFromCSV($path);
    $time_end = microtime(true);

//dividing with 60 will give the execution time in minutes other wise seconds
    $execution_time = ($time_end - $time_start)/60;

//execution time of the script
    echo 'Total Execution Time:'.$execution_time.' Mins' . "\n";
}