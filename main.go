package main

import (
	"fmt"

	"nilphumiphat.assignment.bandprotocol/pkg/transaction"
)

func main() {
	// create new transaction.
	transaction := transaction.NewTransaction("ETH", 50000)

	fmt.Println("Broadcast starting...")

	// broadcast.
	data, err := transaction.Broadcast()

	if err != nil {
		fmt.Println("Broadcast fail")
	} else {
		fmt.Println("Broadcast success is transaction hash => " + data.TxHash)
	}

	// check status
	status, err := transaction.CheckStatus()

	if err != nil {
		fmt.Println("Can not check status please check you internet/network connection")
	} else {
		fmt.Println("Status: " + status.Status)
		fmt.Println("Transaction: " + status.Message)
	}
}
