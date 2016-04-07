package db

import (

	"github.com/jmoiron/sqlx"

	"fmt"
)
var _db *sqlx.DB
// DB error codes (need to formalize this)

const PGDBName = "gwmigration"
const PGDBUserName = "gwdbadmin"
const PGDBPassword = "LigCigUt34"
const PGDBHost = "pdb02-rds.boomlabs.co"
const PGSSLMode = "disable"
const DBMaxConnections = 50
const DBMaxIdleConnections = 5
const DB_SUCCESS = 0
const DB_NO_ROWS = 1
const DB_UNIQUE_VIOLATION = 2

// use for data-migration only
// ==============================================================


const SQL_CREATE_BC_CUSTOMER_TABLE = `
CREATE TABLE IF NOT EXISTS bc_customers (
    bc_customer_id  			INT NOT NULL,
    company						TEXT,
    first_name       			TEXT NOT NULL,
    last_name        			TEXT NOT NULL,
    email		   				TEXT NOT NULL PRIMARY KEY,
    phone		       			TEXT NOT NULL,
    date_created        		TIMESTAMP,
    date_modified       		TIMESTAMP,
    date_imported_db			TIMESTAMP,
    date_csv_created			TIMESTAMP,
    store_credit        		TEXT,
    registration_ip_address     TEXT,
    customer_group_id			INT,
    notes						TEXT,
    tax_exempt_category			TEXT
);`
const SQL_CREATE_BC_ADDRESS_TABLE = `
CREATE TABLE IF NOT EXISTS bc_address (
    id  						INT NOT NULL PRIMARY KEY,
    customer_id					INT NOT NULL,
    first_name       			TEXT NOT NULL,
    last_name        			TEXT NOT NULL,
    company						TEXT,
    street_1		   			TEXT NOT NULL,
    street_2		   			TEXT,
    city		       			TEXT NOT NULL,
    state        				TEXT,
    state_iso2     				TEXT,
    zip_code       				TEXT NOT NULL,
    country						TEXT NOT NULL,
    country_iso2				TEXT,
    phone        				TEXT NOT NULL,
    date_csv_created			TIMESTAMP,
    date_imported_db			TIMESTAMP
);`

const SQL_CREATE_STRIPE_CUSTOMER_TABLE = `
CREATE TABLE IF NOT EXISTS stripe_customers (
	id					TEXT NOT NULL PRIMARY KEY,
	email				TEXT NOT NULL,
	created				TIMESTAMP,
	card_count			INT NOT NULL,
	sub_count			INT NOT NULL,
	date_csv_created	TIMESTAMP,
	date_imported_db	TIMESTAMP
);`

const SQL_DROP_BC_CUSTOMER_TABLE = `
    DROP TABLE bc_customers;`

const SQL_DROP_BC_ADDRESS_TABLE = `
    DROP TABLE bc_address;`

const SQL_DROP_STRIPE_CUSTOMER_TABLE = `
    DROP TABLE stripe_customers;`
/*
test upsert query

const UPSERT_BC_CUSTOMER =
`WITH upsert AS (UPDATE bc_customer SET first_name= :firstname WHERE email=:email RETURNING *)
    INSERT INTO bc_customer (first_name, email) SELECT 'hi', 123D@asdsad.com WHERE NOT EXISTS (SELECT * FROM upsert)`
*/

func GetDBConn() *sqlx.DB {
	if _db != nil {
		return _db
	}

	var err error
	s := fmt.Sprintf("user=%v dbname=%v sslmode=%v password=%v host=%v",
		PGDBUserName,
		PGDBName,
		PGSSLMode,
		PGDBPassword,
		PGDBHost,
	)
	// Note we intentionally DONT use MustConnect since that prints a most unhelpful gross panic
	_db, err = sqlx.Connect("postgres", s)
	if err != nil {
		panic("******** Could not connect to postgres DB. Please make sure postgres is running and check perms")
	}
	// Set the pool connection settings
	_db.SetMaxOpenConns(DBMaxConnections)
	_db.SetMaxIdleConns(DBMaxIdleConnections)
	return _db
}

func CreateBCCustomerTable() {
	sql_smt := SQL_CREATE_BC_CUSTOMER_TABLE
	res, err := GetDBConn().Exec(sql_smt)
	fmt.Println(res, err)
}

func CreateBCAddressTable() {
	sql_smt := SQL_CREATE_BC_ADDRESS_TABLE
	res, err := GetDBConn().Exec(sql_smt)
	fmt.Println(res, err)
}

func CreateStripeCustomerTable() {
	sql_smt := SQL_CREATE_STRIPE_CUSTOMER_TABLE
	res, err := GetDBConn().Exec(sql_smt)
	fmt.Println(res, err)
}

func DropBCCustomerTable() {
	sql_smt := SQL_DROP_BC_CUSTOMER_TABLE
	res, err := GetDBConn().Exec(sql_smt)
	fmt.Println(res, err)
}

func DropBCAddressTable() {
	sql_smt := SQL_DROP_BC_ADDRESS_TABLE
	res, err := GetDBConn().Exec(sql_smt)
	fmt.Println(res, err)
}

func DropStripeCustomerTable() {
	sql_smt := SQL_DROP_STRIPE_CUSTOMER_TABLE
	res, err := GetDBConn().Exec(sql_smt)
	fmt.Println(res, err)
}

func InitSchema() {
	// create all tables that db package needs
	fmt.Println("*******************")
	fmt.Println("Creating tables for BC customers")
	CreateBCCustomerTable()
	fmt.Println("*******************")
	fmt.Println("*******************")
	fmt.Println("Creating tables for BC address")
	CreateBCAddressTable()
	fmt.Println("*******************")
	fmt.Println("*******************")
	fmt.Println("Creating tables for Stripe customers")
	CreateStripeCustomerTable()
	fmt.Println("*******************")
}

// todo: use with caution
func ClearSchema() {
	fmt.Println("*******************")
	fmt.Println("Dropping bc_customers table")
	DropBCCustomerTable()
	fmt.Println("*******************")
	fmt.Println("*******************")
	fmt.Println("Dropping bc_address table")
	DropBCAddressTable()
	fmt.Println("*******************")
	fmt.Println("*******************")
	fmt.Println("Dropping stripe_customer table")
	DropStripeCustomerTable()
	fmt.Println("*******************")
}

