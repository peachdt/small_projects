package db

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"path/filepath"
	"log"
	"encoding/csv"
	"time"
)

/*
	BC ===============================================
 */
const GET_BC_CUSTOMER = `
SELECT * FROM bc_customers WHERE email = '%s';`

const LOAD_BC_CUSTOMER_CSV = `
INSERT INTO bc_customers (
	bc_customer_id,
	first_name,
	last_name,
	phone,
	email,
	date_created,
	date_modified,
	date_imported_db,
	date_csv_created,
	store_credit,
	registration_ip_address,
	customer_group_id,
	notes,
	tax_exempt_category
) VALUES (
	'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s' ,'%s' ,'%s' ,'%s'
);
`
const UPDATE_BC_CUSTOMER_CSV = `
UPDATE bc_customers SET
	bc_customer_id			= %s,
	first_name				= '%s',
	last_name				= '%s',
	phone					= '%s',
	date_created			= '%s',
	date_modified			= '%s',
	date_imported_db		= '%s',
	date_csv_created		= '%s',
	store_credit			= %s,
	registration_ip_address = '%s',
	customer_group_id		= %s,
	notes					= '%s',
	tax_exempt_category		= '%s'
WHERE email = '%s';`

const EXPORT_CUSTOMER_FROM_DB =
`select first_name, last_name, email, date_created from bc_customers;`


/*
	Stripe ===============================================
 */

const GET_STRIPE_CUSTOMER = `
SELECT * FROM stripe_customers WHERE id = '%s';`

const LOAD_STRIPE_CUSTOMER_CSV = `
INSERT INTO stripe_customers (
	id,
	email,
	created,
	card_count,
	sub_count,
	date_csv_created,
	date_imported_db
) VALUES (
	'%s', '%s', '%s', '%s', '%s', '%s', '%s'
);
`
const UPDATE_STRIPE_CUSTOMER_CSV = `
UPDATE stripe_customers SET
	email				= '%s',
	created				= '%s',
	card_count			= %s,
	sub_count			= %s,
	date_csv_created	= '%s',
	date_imported_db	= '%s'
WHERE id = '%s';`

const DELETE_CUSTOMER_FROM_STRIPE = "DELETE FROM stripe_customers WHERE id=$1;"


func UpsertBCCustomerCSV(path string) {
	var sql_smt string
	var index = 1
	var total_customer_count = 0
	var customer_inserted= 0
	var customer_updated = 0
	// find all files under path and save to fileList
	fileList := []string{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	// first item in list is the root, so the actual files start at index 1
	for index < len(fileList)  {
		inFile, _ := os.Open(fileList[index])
		fmt.Println("******************")
		fmt.Println(fmt.Sprintf("Loading %s ...", fileList[index]))
		fmt.Println("******************")
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)
		date_import_db := time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
		for scanner.Scan() {
			total_customer_count += 1
			t := scanner.Text()
			data := strings.Split(t, ",")

			/*
			upsert customer
			 */

			// email is primary key and it's data[4] this case
			result, err := GetDBConn().Exec(fmt.Sprintf(GET_BC_CUSTOMER, data[4]))
			rows_affected, _ := result.RowsAffected()
			if err != nil {
				fmt.Println(err)
			}
			if rows_affected == 0 {
				// do insert
				// only 5 columns to import to db for now
				// todo add more columns if necessary
				sql_smt = fmt.Sprintf(LOAD_BC_CUSTOMER_CSV, data[0], data[1], data[2], data[3], data[4], data[5], data[6], date_import_db , data[12], data[7], data[8], data[9], data[10], data[11])
				_, err := GetDBConn().Exec(sql_smt)
				if err != nil {
					fmt.Println(err)
				} else {
					customer_inserted += 1
				}
			} else {
				// do update
				sql_smt = fmt.Sprintf(UPDATE_BC_CUSTOMER_CSV, data[0], data[1], data[2], data[3], data[5], data[6], date_import_db , data[12], data[7], data[8], data[9], data[10], data[11], data[4])
				_, err := GetDBConn().Exec(sql_smt)
				if err != nil {
					fmt.Println(err)
				} else {
					customer_updated += 1
				}
			}
		}
		index += 1
	}
	fmt.Println("******************")
	fmt.Println("Import completed!")
	fmt.Println(fmt.Sprintf("There are %d new customers inserted and %d customers updated to DB from %d customers in total.", customer_inserted, customer_updated, total_customer_count))
	fmt.Println("******************")
}

// todo: data might be corrupted due to previous whalerock incident
// todo: there are still many duplicate emails with strange states such as one has 1 card 0 sub and the other one has 2 cards 0 sub; or one has 6 cards and others have 3 cards...
// todo: need to do a deep analysis once used on prod
func FilterOutDupsWithoutCard() {
	var id string
	var email string
	var card_count int
	var sub_count int
	var delete_cust_list []string
	var count = 0

	sql_smt := `select id, email, card_count, sub_count from stripe_customers where email in (select email from stripe_customers group by email having count(*) > 1);`
	rows, db_err := GetDBConn().Query(sql_smt)
	if db_err != nil {
		log.Fatal(db_err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &email, &card_count, &sub_count)
		if err != nil {
			log.Fatal(err)
		}
		if card_count == 0 {
			delete_cust_list = append(delete_cust_list, id)
		}
	}
	fmt.Println(delete_cust_list)

	for _, cust := range delete_cust_list {
		fmt.Printf("Deleting customer %s ...\n", cust)
		_, err := GetDBConn().Exec(DELETE_CUSTOMER_FROM_STRIPE, cust)
		if err != nil {
			panic(err)
		}
		count += 1
		fmt.Println(count)
	}
	fmt.Println(fmt.Sprintf("Deleted %d rows.", count))
}

// after taking out dups with no card, leave the most recent created customer and delete all other dups for now
//func GetMostRecentCustInDups() {
//	var id string
//	var email string
//	var created string
//	var temp_list []map[string]string
//	var delete_cust_list []string
//
//	sql_smt := `select id, email, created from stripe_customers where email in (select email from stripe_customers group by email having count(*) > 1);`
//	rows, db_err := GetDBConn().Query(sql_smt)
//	if db_err != nil {
//		log.Fatal(db_err)
//	}
//	defer rows.Close()
//	for rows.Next() {
//		err := rows.Scan(&id, &email, &created)
//		if err != nil {
//			log.Fatal(err)
//		}
//		cust := map[string]string{
//			"id": id,
//			"email": email,
//			"created": created,
//		}
//		temp_list = append(temp_list, cust)
//	}
//	for _, customer := range temp_list {
//
//	}
//}

func ExportCustFromDB() {
	// columns that will be exported from db
	var (
		firstname string
		lastname string
		email_addr string
		created_date time.Time
	)
	var count = 0
	// query database
	rows, db_err := GetDBConn().Query(EXPORT_CUSTOMER_FROM_DB)
	if db_err != nil {
		log.Fatal(db_err)
	}
	defer rows.Close()
	// check if path exists
	path := "../csv/ready/"
	if _, path_err := os.Stat(path); os.IsNotExist(path_err) {
		// if path does not exists, create it
		mkdir_err := os.MkdirAll(path, 0777)
		if mkdir_err != nil {
			panic(fmt.Sprintf("MkdirAll %q: %s", path, mkdir_err))
		}
	}
	// create csv file
	f, create_err := os.Create(path + "_BC_customers.csv")
	fmt.Println(fmt.Sprintf("Saveing DB rows into %s", path + "BC_customers.csv"))
	if create_err != nil {
		fmt.Println(create_err)
	}
	defer f.Close()
	for rows.Next() {
		err := rows.Scan(&firstname, &lastname, &email_addr, &created_date)
		if err != nil {
			log.Fatal(err)
		}
		// parse time.Time to string
		date_split := strings.Split(fmt.Sprintf("%s", created_date), " ")
		date_created := fmt.Sprintf("%s %s", date_split[0], date_split[1])
		// Write data to CSV file
		w := csv.NewWriter(f)
		var record []string
		record = append(record, firstname, lastname, email_addr, date_created)
		w.Write(record)
		w.Flush()
		count += 1
	}
	create_err = rows.Err()
	if create_err != nil {
		log.Fatal(create_err)
	}
	// atomic
	os.Rename("../csv/ready/_BC_customers.csv", "../csv/ready/BC_customers.csv")
	fmt.Println(fmt.Sprintf("There are %d customers saved to csv file.", count))
}

// =================================================================================
func UpsertStripeCustomerCSV(path string) {
	var sql_smt string
	var index = 1
	var total_customer_count = 0
	var customer_inserted= 0
	var customer_updated = 0
	// find all files under path and save to fileList
	fileList := []string{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	// first item in list is the root, so the actual files start at index 1
	for index < len(fileList)  {
		inFile, _ := os.Open(fileList[index])
		fmt.Println("******************")
		fmt.Println(fmt.Sprintf("Loading %s ...", fileList[index]))
		fmt.Println("******************")
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)
		date_import_db := time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
		for scanner.Scan() {
			total_customer_count += 1
			t := scanner.Text()
			data := strings.Split(t, ",")

			/*
			upsert customer
			[id (0), email (1), created (2), card_count (3), sub_count (4), date_csv_created (5)]
			 */

			// id is primary key and it's data[0] this case
			result, err := GetDBConn().Exec(fmt.Sprintf(GET_STRIPE_CUSTOMER, data[0]))
			rows_affected, _ := result.RowsAffected()
			if err != nil {
				fmt.Println(err)
			}
			if rows_affected == 0 {
				// do insert
				// only 5 columns to import to db for now
				// todo add more columns if necessary
				sql_smt = fmt.Sprintf(LOAD_STRIPE_CUSTOMER_CSV, data[0], data[1], data[2], data[3], data[4], data[5], date_import_db)
				_, err := GetDBConn().Exec(sql_smt)
				if err != nil {
					fmt.Println(err)
				} else {
					customer_inserted += 1
				}
			} else {
				// do update
				sql_smt = fmt.Sprintf(UPDATE_STRIPE_CUSTOMER_CSV, data[1], data[2], data[3], data[4], data[5], date_import_db , data[0])
				_, err := GetDBConn().Exec(sql_smt)
				if err != nil {
					fmt.Println(err)
				} else {
					customer_updated += 1
				}
			}
		}
		index += 1
	}
	fmt.Println("******************")
	fmt.Println("Import completed!")
	fmt.Println(fmt.Sprintf("There are %d new customers inserted and %d customers updated to DB from %d customers in total.", customer_inserted, customer_updated, total_customer_count))
	fmt.Println("******************")
}