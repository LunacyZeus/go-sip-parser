package csv_utils

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"sip-parser/pkg/sip"
)

func SaveDataCsv(csvFilePath string, sessions map[string]*sip.SipSession) {
	inFileName := "in_" + csvFilePath
	outFileName := "out_" + csvFilePath

	inFile, err := os.OpenFile(inFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	inCalls := []*PcapCsv{}
	outCalls := []*PcapCsv{}

	for _, session := range sessions {
		if session.Status != sip.COMPLETED {
			continue
		}
		if session.CallBound {
			outCalls = append(outCalls, &PcapCsv{
				CallId:        session.CallID,
				ANI:           session.ANI,
				DNIS:          session.DNIS,
				LRN:           "",
				Via:           session.Via,
				RelatedCallId: session.RelatedCallID,
				OutVia:        session.OutVia,
				InviteTime:    convertTimeStamp(session.InviteTime),
				RingTime:      convertTimeStamp(session.RingTime),
				AnswerTime:    convertTimeStamp(session.AnswerTime),
				HangupTime:    convertTimeStamp(session.HangUpTime),
				Duration:      fmt.Sprintf("%d", session.Duration),
				SrcIP:         session.SrcIP,
				DestIP:        session.DestIP,
			})
		} else {
			inCalls = append(inCalls, &PcapCsv{
				CallId:        session.CallID,
				ANI:           session.ANI,
				DNIS:          session.DNIS,
				LRN:           "",
				Via:           session.Via,
				RelatedCallId: session.RelatedCallID,
				OutVia:        session.OutVia,
				InviteTime:    convertTimeStamp(session.InviteTime),
				RingTime:      convertTimeStamp(session.RingTime),
				AnswerTime:    convertTimeStamp(session.AnswerTime),
				HangupTime:    convertTimeStamp(session.HangUpTime),
				Duration:      fmt.Sprintf("%d", session.Duration),
				SrcIP:         session.SrcIP,
				DestIP:        session.DestIP,
			})
		}
	}

	err = gocsv.MarshalFile(&outCalls, outFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}

	err = gocsv.MarshalFile(&inCalls, inFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}

	fmt.Printf("CSV file->%s and %s created successfully\n", inFileName, outFileName)
}

func AppendDataCsv(csvFilePath string, sessions map[string]*sip.SipSession) {
	// Define the file paths for "in" and "out" CSV files
	inFileName := "in_" + csvFilePath
	outFileName := "out_" + csvFilePath

	// Open the "in" file in append mode
	inFile, err := os.OpenFile(inFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	// Open the "out" file in append mode
	outFile, err := os.OpenFile(outFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Create slices to hold incoming and outgoing call data
	inCalls := []*PcapCsv{}
	outCalls := []*PcapCsv{}

	// Loop over the sessions and categorize them into "in" and "out" calls
	for _, session := range sessions {
		if session.Status != sip.COMPLETED {
			continue
		}
		// Prepare data for outgoing calls
		if session.CallBound {
			outCalls = append(outCalls, &PcapCsv{
				CallId:        session.CallID,
				ANI:           session.ANI,
				DNIS:          session.DNIS,
				LRN:           "",
				Via:           session.Via,
				RelatedCallId: session.RelatedCallID,
				OutVia:        session.OutVia,
				InviteTime:    convertTimeStamp(session.InviteTime),
				RingTime:      convertTimeStamp(session.RingTime),
				AnswerTime:    convertTimeStamp(session.AnswerTime),
				HangupTime:    convertTimeStamp(session.HangUpTime),
				Duration:      fmt.Sprintf("%d", session.Duration),
				SrcIP:         session.SrcIP,
				DestIP:        session.DestIP,
			})
		} else { // Prepare data for incoming calls
			inCalls = append(inCalls, &PcapCsv{
				CallId:        session.CallID,
				ANI:           session.ANI,
				DNIS:          session.DNIS,
				LRN:           "",
				Via:           session.Via,
				RelatedCallId: session.RelatedCallID,
				OutVia:        session.OutVia,
				InviteTime:    convertTimeStamp(session.InviteTime),
				RingTime:      convertTimeStamp(session.RingTime),
				AnswerTime:    convertTimeStamp(session.AnswerTime),
				HangupTime:    convertTimeStamp(session.HangUpTime),
				Duration:      fmt.Sprintf("%d", session.Duration),
				SrcIP:         session.SrcIP,
				DestIP:        session.DestIP,
			})
		}
	}

	// Append the new data to the "out" CSV file
	err = gocsv.MarshalWithoutHeaders(&outCalls, outFile) // Append data to the 'out' file
	if err != nil {
		panic(err)
	}

	// Append the new data to the "in" CSV file
	err = gocsv.MarshalWithoutHeaders(&inCalls, inFile) // Append data to the 'in' file
	if err != nil {
		panic(err)
	}

	// Print success message
	fmt.Printf("CSV data appended to %s and %s successfully\n", inFileName, outFileName)
}
