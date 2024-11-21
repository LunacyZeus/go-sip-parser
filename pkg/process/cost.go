package process

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type CallRecord struct {
	CallID     string
	ANI        string
	DNIS       string
	Via        string
	InviteTime time.Time
	RingTime   time.Time
	AnswerTime time.Time
	HangupTime time.Time
	Duration   int // in milliseconds
	Rate       float64
	RateID     string
	Cost       float64
	Command    string
	Result     string
}

func parseTime(value string) (time.Time, error) {
	layout := "2006-01-02 15:04:05" // Adjust based on your time format
	return time.Parse(layout, value)
}

func CalculateSipCost(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Default delimiter is comma; adjust if needed

	// Read the header
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}
	fmt.Println("Headers:", headers)

	// Read and parse each row
	var records []CallRecord
	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading row:", err)
			continue
		}

		// Parse the row into CallRecord
		inviteTime, _ := parseTime(row[4])
		ringTime, _ := parseTime(row[5])
		answerTime, _ := parseTime(row[6])
		hangupTime, _ := parseTime(row[7])
		duration, _ := strconv.Atoi(row[8])
		rate, _ := strconv.ParseFloat(row[9], 64)
		cost, _ := strconv.ParseFloat(row[11], 64)

		record := CallRecord{
			CallID:     row[0],
			ANI:        row[1],
			DNIS:       row[2],
			Via:        row[3],
			InviteTime: inviteTime,
			RingTime:   ringTime,
			AnswerTime: answerTime,
			HangupTime: hangupTime,
			Duration:   duration,
			Rate:       rate,
			RateID:     row[10],
			Cost:       cost,
			Command:    row[12],
			Result:     row[13],
		}
		records = append(records, record)
	}

	// Output parsed records
	for _, record := range records {
		fmt.Printf("Call Record: %+v\n", record)
	}
}
