package process

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"sip-parser/pkg/utils/rate_utils"
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

var client *telnet.TelnetClient

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

func extractIP(input string) string {
	// 正则表达式匹配IPv4地址
	re := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	// FindString 返回第一个匹配的字符串
	ip := re.FindString(input)
	return ip
}

func handleRow(row []string) (record CallRecord, err error) {
	callerId := row[0]
	ani := row[1]
	dnis := row[2]
	via := row[3]

	// Parse the row into CallRecord
	inviteTime, _ := parseTime(row[4])
	ringTime, _ := parseTime(row[5])
	answerTime, _ := parseTime(row[6])
	hangupTime, _ := parseTime(row[7])

	duration := row[8]
	rate := row[9]
	rateID := row[10]
	cost := row[11]
	command := row[12]
	result := row[13]

	if result != "" || rateID != "" {
		err = fmt.Errorf("calld(%s) already exists", callerId)
		record = CallRecord{
			CallID:     callerId,
			ANI:        ani,
			DNIS:       dnis,
			Via:        via,
			InviteTime: inviteTime,
			RingTime:   ringTime,
			AnswerTime: answerTime,
			HangupTime: hangupTime,
			Duration:   duration,
			Rate:       rate,
			RateID:     rateID,
			Cost:       cost,
			Command:    command,
			Result:     result,
		}
		return
	}

	callerIP := extractIP(row[3])

	aniSip := GetSipPart(row[1])
	dnisSip := GetSipPart(row[2])

	// 建立连接
	err = client.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	// 发送登录命令
	err = client.Login()
	if err != nil {
		err = fmt.Errorf("login->%v", err)
		return
	}

	//log.Println("Login successfully!")

	//fmt.Printf("ani(%s) dnis(%s)\n", row[1], row[2])
	command = fmt.Sprintf("call_simulation %s,5060,%s,%s", callerIP, aniSip, dnisSip)
	log.Printf("[%s] Exec Command-> %s", callerId, command)

	content, err := client.CallSimulation(callerIP, "5060", aniSip, dnisSip)
	if err != nil {
		err = fmt.Errorf("CallSimulation->%v", err)
		return
	}

	_ = client.LoginOut()

	//
	inbound_rate := ""
	inbound_rate_id := ""
	outbound_rate := ""
	outbound_rate_id := ""

	result = ""
	if strings.Contains(content, "No Ingress Resource Found") {
		result = "No Ingress Resource Found"
		log.Printf("[%s]->result: %s", callerId, result)
	} else if strings.Contains(content, "Unauthorized IP Address") {
		result = "Unauthorized IP Address"
		log.Printf("[%s]->result: %s", callerId, result)
	} else if strings.Contains(content, "Ingress Rate Not Found") {
		result = "Ingress Rate Not Found"
		log.Printf("[%s]->result: %s", callerId, result)
	} else {
		inbound_rate, inbound_rate_id, outbound_rate, outbound_rate_id = rate_utils.ParseRateFromContent(callerId, ani, dnis, aniSip, dnisSip, content)
	}

	_ = fmt.Sprintf("%s %s", inbound_rate, inbound_rate_id)

	//fmt.Println(content)

	record = CallRecord{
		CallID:     callerId,
		ANI:        ani,
		DNIS:       dnis,
		Via:        via,
		InviteTime: inviteTime,
		RingTime:   ringTime,
		AnswerTime: answerTime,
		HangupTime: hangupTime,
		Duration:   duration,
		Rate:       outbound_rate,
		RateID:     outbound_rate_id,
		Cost:       cost,
		Command:    command,
		Result:     result,
	}

	return
}

func CalculateSipCost(path string) {
	// 创建客户端实例
	client = telnet.NewTelnetClient("127.0.0.1", "4320")

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Default delimiter is comma; adjust if needed

	// 读取所有行（包括头部）
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	if len(rows) == 0 {
		fmt.Println("CSV file is empty.")
		return
	}

	// 克隆 rows 为 new_rows
	newRows := make([][]string, len(rows))
	copy(newRows, rows)

	// 提取标题行
	headers := rows[0]
	fmt.Println("Headers:", headers)

	// 初始化存储修改后的记录的切片
	var updatedRecords [][]string
	updatedRecords = append(updatedRecords, headers) // 保留标题行

	n := 1
	all_count := len(rows)
	// 实时处理每一行
	for i, row := range rows[1:] { // 跳过标题行
		record, err := handleRow(row)
		if err != nil {
			log.Println("Skip row:", err)
			continue
		}
		fmt.Printf("%d/%d", n, all_count)

		// 修改后的记录写入 new_rows
		newRows[i+1] = recordToRow(record)

		// 实时写入修改后的数据到文件
		err = updateCsv("res_"+path, newRows)
		if err != nil {
			fmt.Println("Error writing CSV:", err)
			return
		}
	}
}

// 将 CallRecord 转换为字符串切片的辅助函数
func recordToRow(record CallRecord) []string {
	return []string{
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
}

// 修改后的 CSV 写入函数
func updateCsv(path string, records [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入所有行
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("Error writing record: %w", err)
		}
	}

	return nil
}
