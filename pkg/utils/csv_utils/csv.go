package csv_utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sip-parser/pkg/sip"
)

func SaveDataCsv(csvFilePath string, sessions map[string]*sip.SipSession) {
	// 写入 CSV 文件的表头
	header := []string{"Call-ID", "ANI", "DNIS", "Invite Time", "Ring Time", "Answer Time", "Hangup Time", "Duration (msec)", "Rate", "Cost"}

	// 创建或打开 CSV 文件
	in_file, err := os.Create("in_" + csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer in_file.Close()

	in_writer := csv.NewWriter(in_file)
	defer in_writer.Flush()

	if err := in_writer.Write(header); err != nil {
		log.Fatal(err)
	}

	// 创建或打开 CSV 文件
	out_file, err := os.Create("out_" + csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out_file.Close()

	out_writer := csv.NewWriter(out_file)
	defer out_writer.Flush()

	if err := out_writer.Write(header); err != nil {
		log.Fatal(err)
	}

	// 写入多个数据行
	for _, session := range sessions {
		if session.Status == sip.COMPLETED { //只解析成功的
			fmt.Println(session.String())
			record := []string{
				session.CallID,
				session.ANI,
				session.DNIS,
				fmt.Sprintf("%d", session.InviteTime),
				fmt.Sprintf("%d", session.RingTime),
				fmt.Sprintf("%d", session.AnswerTime),
				fmt.Sprintf("%d", session.HangUpTime),
				fmt.Sprintf("%d", session.Duration),
				fmt.Sprintf("%d", 0),
				fmt.Sprintf("%d", 0),
			}
			if session.CallBound { //呼出
				if err := out_writer.Write(record); err != nil {
					log.Fatal(err)
				}
			} else {
				if err := in_writer.Write(record); err != nil {
					log.Fatal(err)
				}
			}

		}
	}

	fmt.Printf("CSV file-(%s) created successfully", csvFilePath)
}
