package main

import (
	"BoomPayments/labs/gwmigration/bc"
	"BoomPayments/labs/gwmigration/db"
	"BoomPayments/labs/gwmigration/stripe"

	"flag"
	"fmt"

)

const bc_store = ""
const bc_client = ""
const bc_secret = ""
const bc_store_hash = ""
const bc_csv_path = "../csv/BC/"

func main() {
	action := flag.String("action", "", "main.go action argument")
	path := flag.String("path", "", "main.go action argument")

	flag.Parse()

	switch *action {
	case "export_customers_from_bc":
		fmt.Println(fmt.Sprintf("Exporting Customers from %v . . .", bc_store))
		bc.ExportCustomersFromBC()

	case "export_address_from_bc":
		fmt.Println(fmt.Sprintf("Exporting Addresses from %v . . .", bc_store))
		address_book := bc.ExportAddressFromBC(*path)
		bc.WriteAddrToCSV(address_book)

	case "export_stripe_customer":
		fmt.Println("Exporting Customers from Stripe")
		stripe.ExportCustomersFromStripe()

	case "init_schema":
		fmt.Println("init_shcema")
		db.InitSchema()

	case "clear_schema":
		fmt.Println("clear_schema")
		db.ClearSchema()

	case "load_bc_customer_csv":
		fmt.Println(fmt.Sprintf("Loading BC customer csv from %s", *path))
		db.UpsertBCCustomerCSV(*path)

	case "load_bc_address_csv":
		fmt.Println(fmt.Sprintf("Loading BC address csv from %s", *path))
		db.UpsertBCAddrCSV(*path)

	case "load_stripe_customer_csv":
		fmt.Println(fmt.Sprintf("Loading Stripe customer csv from %s", *path))
		db.UpsertStripeCustomerCSV(*path)

	case "export_bc_customer_from_db":
		db.ExportCustFromDB()

	case "export_bc_address_from_db":
		db.ExportAddrFromDB()

	case "test":
		db.FilterOutDupsWithoutCard()
	default:
		fmt.Println("Arguments missing!")
		fmt.Println("-action=(")
		fmt.Println("    export_customers_from_bc        				Exports customers from BC api to csv")
		fmt.Println("    export_address_from_bc -path       				Exports customer address from BC api to csv, reading customer id from path")
		fmt.Println("    init_schema							Creates all tables in DB")
		fmt.Println("    clear_schema						Drops all tables in DB")
		fmt.Println("    load_bc_customer_csv with -path				Load BC customer csv to DB from path")
		fmt.Println("    load_bc_address_csv with -path				Load BC address csv to DB from path")
		fmt.Println("    export_bc_customer_from_db					Export customers from DB to csv")
		fmt.Println("    export_bc_address_from_db					Export addresses from DB to csv")
		fmt.Println(")")
	}

}
