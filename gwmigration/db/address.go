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

const GET_BC_ADDRESS = `
SELECT * FROM bc_address WHERE id = %s;`

const LOAD_BC_ADDRESS_CSV = `
INSERT INTO bc_address (
	id,
	customer_id,
	first_name,
	last_name,
	company,
	street_1,
	street_2,
	city,
	state,
	state_iso2,
	zip_code,
	country,
	country_iso2,
	phone,
	date_csv_created,
	date_imported_db
) VALUES (
	'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s' ,'%s' ,'%s' ,'%s', '%s', '%s'
);
`

const UPDATE_BC_ADDRESS_CSV = `
UPDATE bc_address SET
	customer_id			= %s,
	first_name			= '%s',
	last_name			= '%s',
	company				= '%s',
	street_1			= '%s',
	street_2			= '%s',
	city				= '%s',
	state				= '%s',
	state_iso2			= '%s',
	zip_code			= '%s',
	country				= '%s',
	country_iso2		= '%s',
	phone				= '%s',
	date_csv_created	= '%s',
	date_imported_db	= '%s'
WHERE id = %s;`

const EXPORT_ADDRESS_FROM_DB =
`select customer_id, first_name, last_name, company, street_1, street_2, city, state, state_iso2, zip_code, country_iso2, phone from bc_address;`


const SELECT_EMAIL_FROM_CUSTOMER_ID =
`select email from bc_customers where bc_customer_id=$1;`

func UpsertBCAddrCSV(path string) {
	var sql_smt string
	var index = 1
	var addresses_inserted = 0
	var addresses_updated = 0
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
			t := scanner.Text()
			data := strings.Split(t, ",")
			/*
				upsert address
			 */
			// id is primary key and it's data[0] in this case
			result, err := GetDBConn().Exec(fmt.Sprintf(GET_BC_ADDRESS, data[0]))
			rows_affected, _ := result.RowsAffected()
			if err != nil {
				fmt.Println(err)
			}
			if rows_affected == 0 {
				// insert
				// only 5 columns to import to db for now
				// todo add more columns if necessary
				sql_smt = fmt.Sprintf(LOAD_BC_ADDRESS_CSV, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], date_import_db)
				_, err := GetDBConn().Exec(sql_smt)
				if err != nil {
					fmt.Println(err)
				} else {
					addresses_inserted += 1
				}
			} else {
				// update
				sql_smt = fmt.Sprintf(UPDATE_BC_ADDRESS_CSV, data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], date_import_db, data[0])
				_, err := GetDBConn().Exec(sql_smt)
				if err != nil {
					fmt.Println(err)
				} else {
					addresses_updated += 1
				}
			}


		}
		index += 1
	}
	fmt.Println("******************")
	fmt.Println("Import completed!")
	fmt.Println(fmt.Sprintf("There are %d new addresses inserted, and %d addresses updated to DB.", addresses_inserted, addresses_updated))
	fmt.Println("******************")
}

func ExportAddrFromDB() {
	// columns that will be exported from db
	var (
		bc_customer_id string
		first_name string
		last_name string
		company string
		street_1 string
		street_2 string
		city string
		zip_code string
		country_id string
		state string
		state_iso2 string
		phone string
		email string
	)
	var count = 0
	// query database
	rows, db_err := GetDBConn().Query(EXPORT_ADDRESS_FROM_DB)
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
	f, create_err := os.Create(path + "_BC_address.csv")
	fmt.Println(fmt.Sprintf("Saveing DB rows into %s", path + "BC_address.csv"))
	if create_err != nil {
		fmt.Println(create_err)
	}
	defer f.Close()
	for rows.Next() {
		/*
		`select customer_id, first_name, last_name, company, street_1, street_2, city, state, zip_code, country, phone from bc_address;`
		 */
		err := rows.Scan(&bc_customer_id, &first_name, &last_name, &company, &street_1, &street_2, &city, &state, &state_iso2, &zip_code, &country_id, &phone)
		if err != nil {
			log.Fatal(err)
		}
		// Write data to CSV file
		w := csv.NewWriter(f)
		var record []string

		// to get email_addr from bc_customers table wieh bc_customer_id.
		// email_addr is used later for php script to use to find customers and update them
		email_field, db_error := GetDBConn().Query(SELECT_EMAIL_FROM_CUSTOMER_ID, bc_customer_id)
		defer email_field.Close()
		for email_field.Next() {
			email_field.Scan(&email)
		}
		if db_error != nil {
			log.Fatal(db_error)
		}

		record = append(record, bc_customer_id, first_name, last_name, company, street_1, street_2, city, state, state_iso2,zip_code, country_id, phone, email)
		w.Write(record)
		w.Flush()
		count += 1
	}
	create_err = rows.Err()
	if create_err != nil {
		log.Fatal(create_err)
	}
	// atomic
	os.Rename("../csv/ready/_BC_address.csv", "../csv/ready/BC_address.csv")
	fmt.Println(fmt.Sprintf("There are %d addresses saved to csv file.", count))
}