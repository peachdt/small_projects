<?php

class OrderGenerator
{
    const CUSTOMER_RANDOM = null;
    
    protected $_shippingMethod = 'freeshipping_freeshipping';
    protected $_paymentMethod = 'cashondelivery';
    
    protected $_customer = self::CUSTOMER_RANDOM;

    protected $_subTotal = 0;
    protected $_order;
    protected $_storeId;

    public function setShippingMethod($methodName)
    {
        $this->_shippingMethod = $methodName;
    }

    public function setPaymentMethod($methodName)
    {
        $this->_paymentMethod = $methodName;
    }
    
    public function setCustomer($customer)
    {
        if ($customer instanceof Mage_Customer_Model_Customer){
            $this->_customer = $customer;
        }
        if (is_numeric($customer)){
            $this->_customer = Mage::getModel('customer/customer')->load($customer);
        }
        else if ($customer === self::CUSTOMER_RANDOM){
            $customers = Mage::getResourceModel('customer/customer_collection');

            $customers
                ->getSelect()
                ->limit(1)
                ->order('RAND()');

            $id = $customers->getFirstItem()->getId();
            
            $this->_customer = Mage::getModel('customer/customer')->load($id);
        }
    }

    public function createOrder($products)
    {
        if (!($this->_customer instanceof Mage_Customer_Model_Customer)){
            $this->setCustomer(self::CUSTOMER_RANDOM);
        }

        $transaction = Mage::getModel('core/resource_transaction');
        $this->_storeId = $this->_customer->getStoreId();
        $reservedOrderId = Mage::getSingleton('eav/config')
            ->getEntityType('order')
            ->fetchNewIncrementId($this->_storeId);

        $currencyCode  = Mage::app()->getBaseCurrencyCode();
        $this->_order = Mage::getModel('sales/order')
            ->setIncrementId($reservedOrderId)
            ->setStoreId($this->_storeId)
            ->setQuoteId(0)
            ->setGlobalCurrencyCode($currencyCode)
            ->setBaseCurrencyCode($currencyCode)
            ->setStoreCurrencyCode($currencyCode)
            ->setOrderCurrencyCode($currencyCode);


        $this->_order->setCustomerEmail($this->_customer->getEmail())
            ->setCustomerFirstname($this->_customer->getFirstname())
            ->setCustomerLastname($this->_customer->getLastname())
            ->setCustomerGroupId($this->_customer->getGroupId())
            ->setCustomerIsGuest(0)
            ->setCustomer($this->_customer);


        $billing = $this->_customer->getDefaultBillingAddress();
        $billingAddress = Mage::getModel('sales/order_address')
            ->setStoreId($this->_storeId)
            ->setAddressType(Mage_Sales_Model_Quote_Address::TYPE_BILLING)
            ->setCustomerId($this->_customer->getId())
            ->setCustomerAddressId($this->_customer->getDefaultBilling())
            ->setCustomerAddress_id($billing->getEntityId())
            ->setPrefix($billing->getPrefix())
            ->setFirstname($billing->getFirstname())
            ->setMiddlename($billing->getMiddlename())
            ->setLastname($billing->getLastname())
            ->setSuffix($billing->getSuffix())
            ->setCompany($billing->getCompany())
            ->setStreet($billing->getStreet())
            ->setCity($billing->getCity())
            ->setCountry_id($billing->getCountryId())
            ->setRegion($billing->getRegion())
            ->setRegion_id($billing->getRegionId())
            ->setPostcode($billing->getPostcode())
            ->setTelephone($billing->getTelephone())
            ->setFax($billing->getFax());
        $this->_order->setBillingAddress($billingAddress);

        $shipping = $this->_customer->getDefaultShippingAddress();
        $shippingAddress = Mage::getModel('sales/order_address')
            ->setStoreId($this->_storeId)
            ->setAddressType(Mage_Sales_Model_Quote_Address::TYPE_SHIPPING)
            ->setCustomerId($this->_customer->getId())
            ->setCustomerAddressId($this->_customer->getDefaultShipping())
            ->setCustomer_address_id($shipping->getEntityId())
            ->setPrefix($shipping->getPrefix())
            ->setFirstname($shipping->getFirstname())
            ->setMiddlename($shipping->getMiddlename())
            ->setLastname($shipping->getLastname())
            ->setSuffix($shipping->getSuffix())
            ->setCompany($shipping->getCompany())
            ->setStreet($shipping->getStreet())
            ->setCity($shipping->getCity())
            ->setCountry_id($shipping->getCountryId())
            ->setRegion($shipping->getRegion())
            ->setRegion_id($shipping->getRegionId())
            ->setPostcode($shipping->getPostcode())
            ->setTelephone($shipping->getTelephone())
            ->setFax($shipping->getFax());

        $this->_order->setShippingAddress($shippingAddress)
            ->setShippingMethod($this->_shippingMethod);

        $orderPayment = Mage::getModel('sales/order_payment')
            ->setStoreId($this->_storeId)
            ->setCustomerPaymentId(0)
            ->setMethod($this->_paymentMethod)
            ->setPoNumber(' – ');

        $this->_order->setPayment($orderPayment);

        $this->_addProducts($products);

        $this->_order->setSubtotal($this->_subTotal)
            ->setBaseSubtotal($this->_subTotal)
            ->setGrandTotal($this->_subTotal)
            ->setBaseGrandTotal($this->_subTotal);

        $transaction->addObject($this->_order);
        $transaction->addCommitCallback(array($this->_order, 'place'));
        $transaction->addCommitCallback(array($this->_order, 'save'));
        $transaction->save();        
    }

    protected function _addProducts($products)
    {
        $this->_subTotal = 0;

        foreach ($products as $productRequest) {
            if ($productRequest['product'] == 'rand') {

                $productsCollection = Mage::getResourceModel('catalog/product_collection');

                $productsCollection->addFieldToFilter('type_id', 'simple');
                Mage::getSingleton('cataloginventory/stock')->addInStockFilterToCollection($productsCollection);

                $productsCollection->getSelect()
                    ->order('RAND()')
                    ->limit(rand($productRequest['min'], $productRequest['max']));

                foreach ($productsCollection as $product){
                    $this->_addProduct(array(
                            'product' => $product->getId(),
                            'qty' => rand(1, 2)
                        ));
                }
            }
            else {
                $this->_addProduct($productRequest);
            }
        }
    }

    protected function _addProduct($requestData)
    {
        $request = new Varien_Object();
        $request->setData($requestData);

        $product = Mage::getModel('catalog/product')->load($request['product']);

        $cartCandidates = $product->getTypeInstance(true)
            ->prepareForCartAdvanced($request, $product);

        if (is_string($cartCandidates)) {
            throw new Exception($cartCandidates);
        }

        if (!is_array($cartCandidates)) {
            $cartCandidates = array($cartCandidates);
        }

        $parentItem = null;
        $errors = array();
        $items = array();
        foreach ($cartCandidates as $candidate) {
            $item = $this->_productToOrderItem($candidate, $candidate->getCartQty());

            $items[] = $item;

            /**
             * As parent item we should always use the item of first added product
             */
            if (!$parentItem) {
                $parentItem = $item;
            }
            if ($parentItem && $candidate->getParentProductId()) {
                $item->setParentItem($parentItem);
            }
            /**
             * We specify qty after we know about parent (for stock)
             */
            $item->setQty($item->getQty() + $candidate->getCartQty());

            // collect errors instead of throwing first one
            if ($item->getHasError()) {
                $message = $item->getMessage();
                if (!in_array($message, $errors)) { // filter duplicate messages
                    $errors[] = $message;
                }
            }
        }
        if (!empty($errors)) {
            Mage::throwException(implode("\n", $errors));
        }

        foreach ($items as $item){
            $this->_order->addItem($item);
        }

        return $items;
    }

    function _productToOrderItem(Mage_Catalog_Model_Product $product, $qty = 1)
    {
        $rowTotal = $product->getFinalPrice() * $qty;

        $options = $product->getCustomOptions();

        $optionsByCode = array();

        foreach ($options as $option)
        {
            $quoteOption = Mage::getModel('sales/quote_item_option')->setData($option->getData())
                ->setProduct($option->getProduct());

            $optionsByCode[$quoteOption->getCode()] = $quoteOption;
        }

        $product->setCustomOptions($optionsByCode);

        $options = $product->getTypeInstance(true)->getOrderOptions($product);

        $orderItem = Mage::getModel('sales/order_item')
            ->setStoreId($this->_storeId)
            ->setQuoteItemId(0)
            ->setQuoteParentItemId(NULL)
            ->setProductId($product->getId())
            ->setProductType($product->getTypeId())
            ->setQtyBackordered(NULL)
            ->setTotalQtyOrdered($product['rqty'])
            ->setQtyOrdered($product['qty'])
            ->setName($product->getName())
            ->setSku($product->getSku())
            ->setPrice($product->getFinalPrice())
            ->setBasePrice($product->getFinalPrice())
            ->setOriginalPrice($product->getFinalPrice())
            ->setRowTotal($rowTotal)
            ->setBaseRowTotal($rowTotal)

            ->setWeeeTaxApplied(serialize(array()))
            ->setBaseWeeeTaxDisposition(0)
            ->setWeeeTaxDisposition(0)
            ->setBaseWeeeTaxRowDisposition(0)
            ->setWeeeTaxRowDisposition(0)
            ->setBaseWeeeTaxAppliedAmount(0)
            ->setBaseWeeeTaxAppliedRowAmount(0)
            ->setWeeeTaxAppliedAmount(0)
            ->setWeeeTaxAppliedRowAmount(0)

            ->setProductOptions($options);

        $this->_subTotal += $rowTotal;

        return $orderItem;
    }
}
