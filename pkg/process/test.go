package process

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"sip-parser/pkg/utils/csv_utils"
)

func TestFunc() {
	clientsFile, err := os.OpenFile("1.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	clients := []*csv_utils.PcapCsv{}

	if err := gocsv.UnmarshalFile(clientsFile, &clients); err != nil { // Load clients from file
		panic(err)
	}
	for index, client := range clients {
		client.Result = "test"
		clients[index] = client
		fmt.Println("Hello", client.CallId)
	}

	if _, err := clientsFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}

	// Save clients to csv file
	if err = gocsv.MarshalFile(&clients, clientsFile); err != nil {
		panic(err)
	}
}
