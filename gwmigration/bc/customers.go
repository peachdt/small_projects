package bc

import (
	"BoomPayments/cs/core_v0/config"
	"BoomPayments/labs/gwmigration/utils"

	"encoding/csv"
	"os"
	"strconv"
	"time"
	"fmt"
	"regexp"
	"log"
)
// Used for exporting customers from data-migration script.
// =============================================================
const API_URL_BASE = "https://api.bigcommerce.com/stores/%s/v2/"
var c = config.GetConfig()
var api_url = fmt.Sprintf(API_URL_BASE, c.BigCommerceApiStoreId)

var http_header = map[string]string{
	"X-Auth-Client": c.BigCommerceApiClientId,
	"X-Auth-Token": c.BigCommerceApiClientToken,
	"Accept": "application/json",
}


type Customer struct {
	Id                  int         `json:"id,omitempty"`
	FirstName           string      `json:"first_name,omitempty"`
	LastName            string      `json:"last_name,omitempty"`
	EmailAddr           string      `json:"email,omitempty"`
	Phone               string      `json:"phone,omitempty"`
	DateCreated         string      `json:"date_created,omitempty"`
	DateModified        string      `json:"date_modified,omitempty"`
	StoreCredit			string		`json:"store_credit,omitempty"`
	RegistrationIPAddr  string		`json:"registration_ip_address,omitempty"`
	CustomerGroupId		int			`json:"customer_group_id,omitempty"`
	Notes				string		`json:"notes,omitempty"`
	TaxExemptCategory	string		`json:"tax_exempt_category"`
}

func ExportCustomersFromBC() {
	// find current timestamp
	time_now := time.Now()
	t := time_now.Format("2006-01-02T15-04-05")
	fmt.Println(fmt.Sprintf("Current timestamp is: %s", t))
	// check if path exists
	path := fmt.Sprintf("../csv/BC/customer/%s/", t)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if path does not exists, create it
		err := os.MkdirAll(path, 0777)
		if err != nil {
			panic(fmt.Sprintf("MkdirAll %q: %s", path, err))
		}
		GetAllCustomers(t)
	}
}

func GetAllCustomers(timestamp string){
	var url string
	var bcr []Customer
	var count = 0
	var sum = 0
	page := 1
	// set url
	url = api_url + fmt.Sprintf("customers?page=%d&limit=250", page)
	fmt.Println("Writing Pgae1.csv ...")
	// call http GET to BC api
	utils.HttpGet(url, http_header, &bcr)
	// write to csv
	count = WriteToCSV(bcr, page, timestamp)
	// update customer sum
	sum = count + sum
	fmt.Println(fmt.Sprintf("... Page1.csv has %d customers.", count))

	for len(bcr) == 250 {
		// more than one page
		page += 1
		// reset url
		url = api_url + fmt.Sprintf("customers?page=%d&limit=250", page)
		fmt.Println(fmt.Sprintf("Writing Pgae%d.csv ...", page))
		utils.HttpGet(url, http_header, &bcr)
		count = WriteToCSV(bcr, page, timestamp)
		fmt.Println(fmt.Sprintf("... Page%d.csv has %d customers.", page, count))
		sum = sum + count
	}
	fmt.Println("Total customers:", sum)
}

func WriteToCSV(input []Customer, page int, timestamp string) int {
	var count = 0
	// Create a csv file
	f, err := os.Create(fmt.Sprintf("../csv/BC/customer/%s/_page%d.csv", timestamp, page))
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	// Write data to CSV file
	w := csv.NewWriter(f)
	for _, obj := range input {
		var record []string
		record = append(record, strconv.Itoa(obj.Id))
		record = append(record, obj.FirstName)
		record = append(record, obj.LastName)
		record = append(record, PhoneVerificationAndCleanUp(obj.Phone))
		record = append(record, obj.EmailAddr)
		record = append(record, ParseTime(obj.DateCreated))
		record = append(record, ParseTime(obj.DateModified))
		record = append(record, obj.StoreCredit)
		record = append(record, obj.RegistrationIPAddr)
		record = append(record, strconv.Itoa(obj.CustomerGroupId))
		record = append(record, obj.Notes)
		record = append(record, obj.TaxExemptCategory)
		date_csv_created := time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
		record = append(record, date_csv_created)

		w.Write(record)
		count += 1
	}
	w.Flush()
	// atomic
	os.Rename(fmt.Sprintf("../csv/BC/customer/%s/_page%d.csv", timestamp, page), fmt.Sprintf("../csv/BC/customer/%s/page%d.csv", timestamp, page))
	return count
}

func PhoneVerificationAndCleanUp(phone string) string {
	if string(phone[0]) == "+" {
		// phone number starts with +, most likely it already has country code in front, so just leave it as it is now
		temp_string := string(phone[1:])
		reg := regexp.MustCompile( "[^0-9]" )
		return "+" + reg.ReplaceAllString(temp_string, "" )
	} else {
		// phone does not start with + means it does not have country code
		reg := regexp.MustCompile( "[^0-9]" )
		temp_phone := reg.ReplaceAllString(phone, "" )
		if len(temp_phone) == 10 {
			// phone number is 10 digits, add +1
			return "+1" + temp_phone
		} else {
			// phone number is invalid, leave as it is
			log.Println(fmt.Sprintf("Phone number: %s is not 10 digits.", temp_phone))
			return temp_phone
		}
	}
}

func ParseTime(datetime string) string {
	date_parsed, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", datetime)
	if err != nil {
		panic(err)
	}
	return date_parsed.Format("2006-01-02 15:04:05.999999 +0000")
}

