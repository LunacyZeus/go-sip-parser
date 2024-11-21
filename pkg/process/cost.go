package process

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"sip-parser/pkg/utils/rate"
	"sip-parser/pkg/utils/telnet"
	"strings"
)

type CallRecord struct {
	CallID     string
	ANI        string
	DNIS       string
	Via        string
	InviteTime string
	RingTime   string
	AnswerTime string
	HangupTime string
	Duration   string // in milliseconds
	Rate       string
	RateID     string
	Cost       string
	Command    string
	Result     string
}

func parseTime(value string) (string, error) {
	//layout := "2006-01-02 15:04:05" // Adjust based on your time format
	return value, nil
}

func GetCallerIP(input string) string {
	// 正则表达式：匹配@后面的IP地址
	re := regexp.MustCompile(`@(\d+\.\d+\.\d+\.\d+)`)
	// 查找所有匹配
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func GetSipPart(input string) string {
	// 去除 <sip: 和 @ 后面的内容
	if strings.Contains(input, "<sip:") {
		// 截取 <sip: 后的部分，去除后面的 @ 和 IP
		start := strings.Index(input, "<sip:") + len("<sip:")
		end := strings.Index(input[start:], "@")
		if end != -1 {
			// 打印号码部分
			return input[start : start+end]
		}
	} else {
		// 如果不是 <sip: 格式，按 @ 分割并打印号码部分
		parts := strings.Split(input, "@")
		if len(parts) > 0 {
			return parts[0] // 打印号码部分
		}
	}
	return ""
}

func writeCsv(csvPath string, headers []string, records []CallRecord) {
	// Write the modified records to a new CSV file
	outputFile, err := os.Create(csvPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Write the header to the new CSV file
	if err := writer.Write(headers); err != nil {
		fmt.Println("Error writing header:", err)
		return
	}

	// Write each modified record to the new CSV file
	for _, record := range records {
		row := []string{
			record.CallID,
			record.ANI,
			record.DNIS,
			record.Via,
			record.InviteTime,
			record.RingTime,
			record.AnswerTime,
			record.HangupTime,
			record.Duration,
			record.Rate,
			record.RateID,
			record.Cost,
			record.Command,
			record.Result,
		}

		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing record:", err)
		}
	}

	fmt.Printf("Modified CSV written successfully to '%s'\n", csvPath)
}

func CalculateSipCost(path string) {
	// 创建客户端实例
	client := telnet.NewTelnetClient("127.0.0.1", "4320")

	// 建立连接
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	// 发送登录命令
	err = client.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer client.LoginOut()

	log.Println("Login successfully!")

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

		callerIP := GetCallerIP(row[2])

		ani := GetSipPart(row[1])
		dnis := GetSipPart(row[2])

		//fmt.Printf("ani(%s) dnis(%s)\n", row[1], row[2])
		command := fmt.Sprintf("call_simulation %s,5060,%s,%s", callerIP, ani, dnis)
		log.Printf("[%s] Exec Command-> %s", row[0], command)

		content, err := client.CallSimulation(callerIP, "5060", ani, dnis)
		if err != nil {
			log.Println("CallSimulation", err)
			return
		}

		result := ""
		if strings.Contains(content, "No Ingress Resource Found") {
			result = "No Ingress Resource Found"
			log.Printf("callId(%s)->%s", row[0], result)
		} else if strings.Contains(content, "Unauthorized IP Address") {
			result = "Unauthorized IP Address"
			log.Printf("callId(%s)->%s", row[0], result)
		} else {
			rate.ParseRateFromContent(content)
		}

		//fmt.Println(content)

		record := CallRecord{
			CallID:     row[0],
			ANI:        row[1],
			DNIS:       row[2],
			Via:        row[3],
			InviteTime: inviteTime,
			RingTime:   ringTime,
			AnswerTime: answerTime,
			HangupTime: hangupTime,
			Duration:   row[8],
			Rate:       row[9],
			RateID:     row[10],
			Cost:       row[11],
			Command:    command,
			Result:     result,
		}
		records = append(records, record)
	}

	// Output parsed records
	/*
		for _, record := range records {
			//fmt.Printf("Call Record: %+v\n", record)
		}
	*/
	writeCsv("out.csv", headers, records)
}
