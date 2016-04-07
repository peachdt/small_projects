package stripe

import (

	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go"
	"fmt"
	"time"
	"os"
	"encoding/csv"
	"strconv"
	"strings"
)

const url = "https://api.stripe.com/v1/customers?limit=100"
// stripe prod key
//const api_key = "sk_live_D5GuNacp8hfX76tg3cOt1Aha"

// stripe test key
const api_key = "sk_test_MHcGObS6pzeFxEJvfePx7tg7"

func ExportCustomersFromStripe() {
	// find current timestamp
	time_now := time.Now()
	t := time_now.Format("2006-01-02T15-04-05")
	fmt.Println(fmt.Sprintf("Current timestamp is: %s", t))
	// check if path exists
	path := fmt.Sprintf("../csv/stripe/customer/%s/", t)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if path does not exists, create it
		err := os.MkdirAll(path, 0777)
		if err != nil {
			panic(fmt.Sprintf("MkdirAll %q: %s", path, err))
		}
		GetAllStripeCustomers(t)
	}
}

func GetAllStripeCustomers(t string) {
	var customer_list []stripe.Customer
	stripe.Key = api_key

	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("limit", "", "100")
	i := customer.List(params)
	for i.Next() {
		c := i.Customer()
		customer_list = append(customer_list, *c)
	}
	WriteToCSV(customer_list, t)
}

func WriteToCSV(customer []stripe.Customer, timestamp string) int {
	var count = 0
	// Create a csv file
	f, err := os.Create(fmt.Sprintf("../csv/stripe/customer/%s/_stripe_customers.csv", timestamp))
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	// Write data to CSV file
	w := csv.NewWriter(f)

	for _, obj := range customer {
		var record []string
		record = append(record, obj.ID)
		record = append(record, obj.Email)
		created_time := ParseUnixTime(int(obj.Created))
		record = append(record, created_time)
		record = append(record, fmt.Sprint(obj.Sources.Count))
		record = append(record, fmt.Sprint(obj.Subs.Count))
		date_csv_created := time.Now().UTC().Format("2006-01-02 15:04:05 +0000")
		record = append(record, date_csv_created)
		w.Write(record)
		count += 1
	}
	w.Flush()
	// atomic
	os.Rename(fmt.Sprintf("../csv/stripe/customer/%s/_stripe_customers.csv", timestamp), fmt.Sprintf("../csv/stripe/customer/%s/stripe_customers.csv", timestamp))
	fmt.Println(fmt.Sprintf("There are %d customers exported from Stripe API", count))
	return count
}

func ParseUnixTime(ut int) string {
	var time_after_parse string
	i, err := strconv.ParseInt(fmt.Sprint(ut), 10, 64)
	if err != nil {
		panic(err)
	}
	time := time.Unix(i, 0)
	time_string := fmt.Sprint(time.UTC().Format("2006-01-02 15:04:05 +0000"))
	time_parse := strings.Split(time_string, " ")
	time_after_parse = fmt.Sprintf("%s %s %s", time_parse[0], time_parse[1], time_parse[2])
	return time_after_parse
}

