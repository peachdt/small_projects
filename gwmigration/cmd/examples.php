<?php

require_once 'app/Mage.php';
umask(0);
Mage::app('default');

require_once 'OrderGenerator.php';

//function create_order_with_custom_products()
//{
//    $orderGenerator = new OrderGenerator();
//
//    $orderGenerator->createOrder(array(
//        // Add configurable product
//        array(
//            'product' => 418,
//            'super_attribute' => array(
//                92 => 26,
//                180 => 79
//            ),
//            'qty' => 1
//        ),
//        // Add products with custom options
//        array(
//            'product' => 553,
//            'options' => array(
//                11 => 8
//            ),
//            'qty' => 2
//        ),
//        array(
//            'product' => 553,
//            'options' => array(
//                11 => 9
//            ),
//            'qty' => 1
//        ),
//        // Add 1-3 random simple products
//        array(
//               'product' => 'rand',
//               'min' => 1,
//               'max' => 3
//        ),
//    ));
//}

function create_random_orders($qty)
{
    $orderGenerator = new OrderGenerator();

    for ($i = 0; $i < $qty; $i++)
    {
        $orderGenerator->setCustomer(OrderGenerator::CUSTOMER_RANDOM);

        $orderGenerator->createOrder(array(
            array(
                'product' => 'rand',
                'min' => 1,
                'max' => 3
        )));
    }
}

//function create_customers_with_order($qty)
//{
//    $orderGenerator = new OrderGenerator();
//    $customerGenerator = new CustomerGenerator();
//
//    for ($i = 0; $i < $qty; $i++)
//    {
//        $customer = $customerGenerator->createCustomer(array(
//                'account' => array(
//                    'lastname' => 'Lastname' . rand(),
//                )
//            ));
//
//        $orderGenerator->setCustomer($customer);
//
//        $orderGenerator->createOrder(array(
//            array(
//            'product' => 'rand',
//            'min' => 1,
//            'max' => 3
//        )));
//    }
//}
//
//
//function create_customers_and_random_orders($customersQty, $ordersQty)
//{
//    $orderGenerator = new OrderGenerator();
//    $customerGenerator = new CustomerGenerator();
//
//    $customers = array();
//
//    for ($i = 0; $i < $customersQty; $i++)
//    {
//        $customers []= $customerGenerator->createCustomer(array(
//            'address' => array(
//                '_item1' => array('lastname' => 'Lastname{id}'),
//            )
//        ));
//    }
//
//
//    for ($i = 0; $i < $ordersQty; $i++){
//
//        $customer = $customers[rand(0, $customersQty - 1)];
//        $orderGenerator->setCustomer($customer);
//
//        $orderGenerator->createOrder(array(
//            array(
//                'product' => 'rand',
//                'min' => 1,
//                'max' => 3
//        )));
//    }
//}
