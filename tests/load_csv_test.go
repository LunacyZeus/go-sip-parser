package tests

import (
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"sip-parser/pkg/sip"
	"sip-parser/pkg/utils"
	"sip-parser/pkg/utils/csv_utils"
	"testing"
)

func loadCsvFile(path string) []*csv_utils.PcapCsv {
	csvFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	rows := []*csv_utils.PcapCsv{}

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

	return rows
}

func getCommand(row *csv_utils.PcapCsv) string {
	var callerIP string

	if row.Via == "" && row.SrcIP != "" {
		callerIP = utils.ExtractIP(row.SrcIP)

	} else {
		callerIP = utils.ExtractIP(row.Via)
	}

	aniSip := sip.GetSipPart(row.ANI)
	dnisSip := sip.GetSipPart(row.DNIS)
	callerPort := "5060"

	// 构建命令
	command := fmt.Sprintf("call_simulation %s,%s,%s,%s\r\n", callerIP, callerPort, aniSip, dnisSip)
	return command
}

func TestLoadCsv(t *testing.T) {
	csv1 := loadCsvFile("../data/test_csv_1.csv")
	csv2 := loadCsvFile("../data/test_csv_2.csv")

	csv1_row1 := csv1[0]
	csv2_row1 := csv2[0]

	command1 := getCommand(csv1_row1)
	command2 := getCommand(csv2_row1)

	t.Logf("command1: %s", command1)
	t.Logf("command2: %s", command2)
}
