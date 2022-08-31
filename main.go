// hotel_billing project main.go
package main

import (
	"fmt"
	"hotel_billing/shredder"
	"log"
)

func main() {
	fmt.Println("Init hotel_billing!")

	err := shredder.DistributeShred()
	if err == nil {
		fmt.Println("Ok execute")
	} else {
		fmt.Println("Bag execute")
		log.Fatal(err)
	}

}
