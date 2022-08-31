// billing project hotel_billing.go
package billing

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type billingHotel struct {
	customerName string
	total        int64
	vat          int64
	vatTotal     int64
}

func HelloBilling() {
	fmt.Println("Hello Billing!")
}
func ProcessBilling(summary chan<- string, nameFile string, pathSummary string, processN int) {
	mapBilling, err := loadBilling(nameFile)
	if err == nil {
		tempSum, err := writeBillingSummary(mapBilling, pathSummary, processN)
		if err == nil {
			//	fmt.Printf("\nThe file summary is %s, file is %s", tempSum, nameFile)
			summary <- tempSum //

		}
	} else {
		log.Fatalf("The operation has faild because the next cause, %s", err)
	}

}
func writeBillingSummary(mapBills map[string]billingHotel, pathSummary string, processN int) (string, error) {
	var nameFile string
	var fTotal, fVat, vatTotal float64
	now := time.Now()
	nameFile = fmt.Sprintf("%s%sbillingSummary.-%d-%s.csv", pathSummary, string(os.PathSeparator), processN, strings.Replace(now.Format(time.RFC3339), ":", "", 3))
	f, err := os.Create(nameFile)
	defer f.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(f)
	defer w.Flush()
	fTotal = 0
	fVat = 0
	vatTotal = 0
	totalRecord := 0
	for _, v := range mapBills {
		fTotal += (float64)(v.total) / 100.00
		fVat += (float64)(v.vat) / 100.00
		vatTotal += (float64)(v.vatTotal) / 100.00
		totalRecord++
	}

	records := [][]string{{"TOTAL", "VAT", "VAT_TOTAL", "TOTAL_RECORD"},
		{fmt.Sprintf("%4.2f", fTotal),
			fmt.Sprintf("%4.2f", fVat),
			fmt.Sprintf("%4.2f", vatTotal),
			fmt.Sprintf("%d", totalRecord)}}

	if err := w.WriteAll(records); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	return nameFile, nil
}
func loadBilling(nameFile string) (map[string]billingHotel, error) {
	mapBilling := make(map[string]billingHotel)
	file, err := os.Open(nameFile)
	if err != nil {
		log.Fatalf("impossible to open file %s", err)
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		arr, err := strToInt64(record)
		if err != nil {
			return nil, err
		}
		bill := billingHotel{
			customerName: record[1],
			total:        arr[0],
			vat:          arr[1],
			vatTotal:     arr[2]}
		mapBilling[record[0]] = bill
	}
	return mapBilling, nil
}

func strToInt64(records []string) ([]int64, error) {
	arrInt := make([]int64, 3)
	var c int
	c = 0
	for i, v := range records {
		if i > 1 {
			if s, err := strconv.ParseInt(v, 10, 64); err == nil {
				arrInt[c] = s
			} else {
				return nil, err
			}
			c++
		}
	}
	return arrInt, nil
}
