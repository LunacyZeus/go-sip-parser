package process

/*
func HandleRowTest(row []string) (record CallRecord, err error) {
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

	if result != "" || rate != "" || rateID != "" || cost != "" {
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

	cost = "12345"

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

func CalculateSipCostTest(path string) {
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

	// 实时处理每一行
	for i, row := range rows[1:] { // 跳过标题行
		record, err := HandleRowTest(row)
		if err != nil {
			log.Println("Skip row:", err)
			continue
		}

		// 将修改后的数据记录转回字符串切片格式
		//fmt.Println(record)

		// 修改后的记录写入 new_rows
		newRows[i+1] = recordToRow(record)

		// 实时写入修改后的数据到文件
		err = updateCsv(path, newRows)
		if err != nil {
			fmt.Println("Error writing CSV:", err)
			return
		}
	}
}
*/
