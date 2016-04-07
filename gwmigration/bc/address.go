package bc

import (
//	"BoomPayments/cs/core_v0/config"
	"BoomPayments/labs/gwmigration/utils"

	"encoding/csv"
	"os"
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
	"log"
	"time"
	"strconv"
)

type Address struct {
	Id					int					`json:"id"`
	CustomerId			int					`json:"customer_id"`
	FirstName			string				`json:"first_name"`
	LastName			string				`json:"last_name"`
	Company				string				`json:"company,omitempty"`
	Street1				string				`json:"street_1"`
	Street2				string				`json:"street_2"`
	City				string				`json:"city"`
	State				string				`json:"state"`
	ZipCode				string				`json:"zip"`
	Country				string				`json:"country"`
	CountryISO2			string				`json:"country_iso2"`
	Phone				string				`json:"phone"`
}

var states_dict = map[string]string {
	"Alabama": "AL",
	"Alaska": "AK",
	"American Samoa": "AS",
	"Arizona": "AZ",
	"Arkansas": "AR",
	"California": "CA",
	"Colorado": "CO",
	"Connecticut": "CT",
	"Delaware": "DE",
	"District Of Columbia": "DC",
	"Federated States Of Micronesia": "FM",
	"Florida": "FL",
	"Georgia": "GA",
	"Guam": "GU",
	"Hawaii": "HI",
	"Idaho": "ID",
	"Illinois": "IL",
	"Indiana": "IN",
	"Iowa": "IA",
	"Kansas": "KS",
	"Kentucky": "KY",
	"Louisiana": "LA",
	"Maine": "ME",
	"Marshall Islands": "MH",
	"Maryland": "MD",
	"Massachusetts": "MA",
	"Michigan": "MI",
	"Minnesota": "MN",
	"Mississippi": "MS",
	"Missouri": "MO",
	"Montana": "MT",
	"Nebraska": "NE",
	"Nevada": "NV",
	"New Hampshire": "NH",
	"New Jersey": "NJ",
	"New Mexico": "NM",
	"New York": "NY",
	"North Carolina": "NC",
	"North Dakota": "ND",
	"Northern Mariana Islands": "MP",
	"Ohio": "OH",
	"Oklahoma": "OK",
	"Oregon": "OR",
	"Palau": "PW",
	"Pennsylvania": "PA",
	"Puerto Rico": "PR",
	"Rhode Island": "RI",
	"South Carolina": "SC",
	"South Dakota": "SD",
	"Tennessee": "TN",
	"Texas": "TX",
	"Utah": "UT",
	"Vermont": "VT",
	"Virgin Islands": "VI",
	"Virginia": "VA",
	"Washington": "WA",
	"West Virginia": "WV",
	"Wisconsin": "WI",
	"Wyoming": "WY",
}

func GetCustomerIDFromCSV(path string) []string{
	fmt.Println("******************")
	fmt.Println("First, load customer id from csv ...")
	fmt.Println("******************")
	var file_index = 1
	var customer_id_list  []string
	// to get address, we need customer_id first
	//get customer_id from csv file
	fileList := []string{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if len(fileList) < 2 {
		// this means there are no csv file under this path
		log.Println(fmt.Sprintf("There are no csvs to read from %s", path))
		return nil
	} else {
		for file_index < len(fileList) {
			inFile, _ := os.Open(fileList[file_index])
			fmt.Println("******************")
			fmt.Println(fmt.Sprintf("Loading %s ...", fileList[file_index]))
			fmt.Println("******************")
			defer inFile.Close()
			scanner := bufio.NewScanner(inFile)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				t := scanner.Text()
				data := strings.Split(t, ",")
				// only 5 columns to import to db for now
				// todo add more columns if necessary
				customer_id_list = append(customer_id_list, data[0])
			}
			file_index += 1
		}
		return customer_id_list
	}
}


func ExportAddressFromBC(path string) []Address{
	var url string
	var address_book []Address
	page := 1
	customer_id_list := GetCustomerIDFromCSV(path)
	if customer_id_list != nil {
		for _, id := range customer_id_list {
			var addr []Address
			url = api_url + fmt.Sprintf("customers/%s/addresses?page=%d&limit=250", id, page)
			// call http GET to BC api
			utils.HttpGet(url, http_header, &addr)
			if addr != nil {
				address := addr[len(addr)-1]
				address.Phone = PhoneVerificationAndCleanUp(address.Phone)
				address_book = append(address_book, address)
			}
		}
	}
	return address_book
}

func WriteAddrToCSV(address_book []Address) int {
	var count = 0
	// create path if does not exist
	time_now := time.Now()
	t := time_now.Format("2006-01-02T15-04-05")
	fmt.Println(fmt.Sprintf("Current timestamp is: %s", t))
	// check if path exists
	path := fmt.Sprintf("../csv/BC/address/%s/", t)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if path does not exists, create it
		err := os.MkdirAll(path, 0777)
		if err != nil {
			panic(fmt.Sprintf("MkdirAll %q: %s", path, err))
		}
	}
	f, err := os.Create(fmt.Sprintf("../csv/BC/address/%s/_BC_customer_address.csv", t))
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	// Write data to CSV file
	w := csv.NewWriter(f)
	for _, address := range address_book {
		var record []string
		record = append(record, strconv.Itoa(address.Id))
		record = append(record, strconv.Itoa(address.CustomerId))
		record = append(record, address.FirstName)
		record = append(record, address.LastName)
		record = append(record, address.Company)
		record = append(record, address.Street1)
		record = append(record, address.Street2)
		record = append(record, address.City)
		record = append(record, address.State)
		record = append(record, StateHash(address.State))
		record = append(record, address.ZipCode)
		record = append(record, address.Country)
		record = append(record, address.CountryISO2)
		record = append(record, address.Phone)
		date_csv_created := time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
		record = append(record, date_csv_created)

		w.Write(record)
		count += 1
	}
	w.Flush()
	os.Rename(fmt.Sprintf("../csv/BC/address/%s/_BC_customer_address.csv", t), fmt.Sprintf("../csv/BC/address/%s/BC_customer_address.csv", t))
	fmt.Println(fmt.Sprintf("Exported %d addresses from BC api to address csv.", count))
	return count
}

func StateHash(state string) string {
	// used for php script to find region id to import addresses in Magento
	if val, ok := states_dict[state]; ok {
		return val
	}
	return ""
}