# GW Data Migration
To migrate customer data from BigCommerce and Stripe to Magento.

## Flow Chart:
![Alt text](https://github.com/BoomPayments/labs/blob/master/gwmigration/utils/flow_chart.png)

## Where to start
Access all commands from `cmd/main.go` by running `go run main.go` under `cmd` directory.
From there, it lists all the possible arguements within main.go.
As of now, argument list looks like:

| Actions                  | Description                                       | Second actions |
| ------------------------ |:-------------------------------------------------:|:----------------:|
| init_schema              | creates all tables in RDS                         | N/A|
| clear_schema             | drops all tables in RDS                           | N/A|
|export_customers_from_bc  |use BC api to export customer data to local csv    | N/A|
|export_address_from_bc    |use BC api to export address data to local csv     | -path |
|export_stripe_customer    |use Stripe api to export customer data to local csv| N/A |
|load_bc_customer_csv      |import bc customer csv to RDS                      | -path |
|load_bc_address_csv       |import bc address csv to RDS                       | -path |
|load_stripe_customer_csv  |import stripe customer csv to RDS                  | -path |
|export_bc_customer_from_db|export bc customer data from RDS to local csv      | N/A |
|export_bc_address_from_db |export bc address data from RDS to local csv       | N/A |

## Steps
* Below assumes all commands run under cmd directory
* **Note:** Please modify psql creds to use data-migration RDS DB creds. Also modify BC creds as well in base.go

* Create all tables in RDS:
```
go run main.go -action=init_schema
```
Three tables will be created if not exists in RDS: `bc_customers, bc_address and stripe_customers`

* To drop all tables in RDS: 
```
go run main.go -action=clear_schema
```
There is no checking for prod env therefore be cautious while using `clear_schema`.

### BigCommerce
Access customer and address data through BigCommerce API

[BigCommerce API Doc](https://developer.bigcommerce.com/api/stores/v2/customers)
#### Customer
* First, get all customers from BigCommerce with all fields provided from BC API.
```
go run main.go -action=export_customers_from_bc
```
BigCommerce API only allows max 250 customers to be displayed in one API call, therefore the results are saved in multiple CSV files based on the actual number of customers in BC store.

CSV files are named as page1.csv, page2.csv ... etc. They will be located in `gwmigration/csv/BC/customer/<timestamp>/` directory. `<timestamp>` will have the format as `2015-11-11T12-00-00` in local timezone.

Phone number will be parsed to take out all the non-numerical characters coming from BC API. If the phone number has 10 digits after parsing, a `+1` will be added in the front as default country code.

> Email clean up is under development.

A new field `date_csv_created` is added in the CSV as the time when CSV file is created. It is in the format of `2015-11-11 12:00:00 +0000` in UTC. The other two fields `date_created` and `date_modified` from BC API will follow the same format in UTC.

#### Address
* To get the addresses for the exported customers, run 
```
go run main.go -action=export_address_from_bc -path="/path/to/customer/csv"
```
This command will take all the customer ids from customer CSV files exported from previous step to query BC API.

Path `/path/to/customer/csv` will only need to be the path to the `<timestamp>` directory, for example `gwmigraion/csv/BC/customer/2015-11-11T12-00-00/` since the script will detect all the individual files under this directory.

Result will be created as one CSV file under `gwmigraion/csv/BC/address/<timestamp>/BC_customer_address.csv`.

Since Magento uses Region Id to find state in adress by using , the script will parse the full name of state such as `California` to `CA` and use ISO 3166-1 alpha-2 country code for country such as `US`.

`date_csv_created` is also added in this CSV with the same format as in previous customer step.

### Stripe
#### Customer
* To get customer data from Stripe API, run:
```
go run main.go -action=export_stripe_customer
```
API key is set to use test Stripe account for now. Time from Stripe API is in Unix Time, and the script will parse it to use standard UTC format. 

Result will be saved in `gwmigration/csv/stripe/customer/<timestamp>/stripe_customers.csv`.

* Process to filter out customers with duplicate email is under development.

### RDS
#### Import to DB
* To load BC customer csv file to DB, run
```
go run main.go -action=load_bc_customer_csv -path="/path/to/customer/csv"
```
Path also only needs to be at `<timestamp>` directory.

If customer is not found in DB based on `eamil`, it will be added as a new row. If customer is found in DB, it will be updated under `bc_customers` table.

* To load BC address csv file to DB, run
```
go run main.go -action=load_bc_address_csv -path="/path/to/address/csv"
```
Path also only needs to be at `<timestamp>` directory.

Same as customer, script will upsert address in DB under `bc_address` table.

* To load Stripe customer csv file to DB, run
```
go run main.go -action=load_stripe_customer_csv -path="/path/to/stripe/customer/csv"
```
Path also only needs to be at `<timestamp>` directory.
Same as previous, script will upsert stripe customer in DB under `stripe_customers` table


#### Export from DB
* To export BC customer data from DB to CSV file, run:
```
go run main.go -action=export_bc_customer_from_db
```
Result will be created under `gwmigration/csv/ready/BC_customers.csv`. This CSV file is ready to upload to Magento.

* To export BC address data from DB to CSV file, run:
```
go run main.go -action=export_bc_address_from_db
```
Result will be created under `gwmigration/csv/ready/BC_address.csv`. This CSV file is ready to upload to Magento.

> Export Stripe customer is still under development.

### Upload to Magento
**Note:** Please set timezone to UTC from Magento Admin Panel before importing anything.
* PHP scripts will be used to upload CSV files to Magento on remote Magento box as root user.

`importCustomer.php` is used to upload BC customer data, and `importAddress.php` is used to upload BC address data.
PHP script uses Mage.php directly from Magento server to call native functions in Models. 

* To use php scripts, `scp` both php files from `gwmigration/cmd/` to remote Magento box. Also `scp` CSV files under `gwmigration/csv/ready/` to `/var/www/magento/var/import/` on remote Magento box. 
* Run `php importCustomer.php` to upload BC customer data. Script will look inside `.../var/import/` to find `BC_customers.csv` and start loading customers.
* Run `php importAddress.php` to upload BC address data. Script will find `BC_address.csv` and start loading addresses.
