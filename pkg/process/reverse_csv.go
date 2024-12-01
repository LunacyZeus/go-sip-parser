package process

import (
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"io"
	"log"
	"os"
	"path/filepath"
	"sip-parser/pkg/utils/csv_utils"
)

func ReverseCsv(path string) {
	csvFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	rows := []*csv_utils.PcapCsv{}
	new_rows := []*csv_utils.PcapCsv{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		//r.LazyQuotes = true
		//r.Comma = '.'
		r.FieldsPerRecord = -1
		return r // Allows use dot as delimiter and use quotes in CSV
	})

	if err := gocsv.UnmarshalFile(csvFile, &rows); err != nil { // Load clients from file
		panic(err)
	}

	n := 1
	//all_count := len(rows)

	for index, _ := range rows {
		row := rows[index]
		new_row, err := csv_utils.ReverseRow(row)
		if err != nil {
			log.Println("Skip row:", err)
			continue
		}
		//log.Printf("processing->%d/%d", n, all_count)

		new_rows = append(new_rows, new_row)
		n += 1
	}

	fileName := filepath.Base(path)
	fileName = "converted_" + fileName

	csvWriteFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	//每操作一次写入一次
	err = gocsv.MarshalFile(&new_rows, csvWriteFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}

	csvWriteFile.Close()
}
