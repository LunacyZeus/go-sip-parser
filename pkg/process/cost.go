package process

import (
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"sip-parser/pkg/sip"
	"sip-parser/pkg/utils"
	"sip-parser/pkg/utils/csv_utils"
	"sip-parser/pkg/utils/telnet"
	"strings"
	"sync"
)

type CallRecord struct {
	CallID        string
	ANI           string
	DNIS          string
	Via           string
	RelatedCallID string
	OutVia        string
	InviteTime    string
	RingTime      string
	AnswerTime    string
	HangupTime    string
	Duration      string // in milliseconds
	InRate        string
	InRateID      string
	InCost        string
	OutRate       string
	OutRateID     string
	OutCost       string
	Command       string
	Result        string
}

func handleRow(pool *telnet.TelnetClientPool, row *csv_utils.PcapCsv) (err error) {
	// 创建客户端实例
	//client := telnet.NewTelnetClient("127.0.0.1", "4320")
	// 获取一个客户端实例

	callerId := row.CallId
	//via := row.Via

	//relatedCallID := row.RelatedCallId
	outVia := row.OutVia

	// Parse the row into CallRecord
	//inviteTime := row.InviteTime
	//ringTime := row.RingTime
	//answerTime := row.AnswerTime
	//hangupTime := row.HangupTime

	//duration := row.Duration
	//inRate := row.InRate
	inRateID := row.InRateID
	//inCost := row.InCost
	//outRate := row.OutRate
	//outRateID := row.OutRateID
	//outCost := row.OutCost
	inTrunkId := row.InTrunkId
	outTrunkId := row.OutTrunkId

	command := row.Command
	result := row.Result

	client, err := pool.Get("127.0.0.1", "4320")
	if err != nil {
		log.Println(err)
		return
	}

	if !client.IsAuthentication {
		// 发送登录命令
		err = client.Login()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("[%s] Successfully logged in!", client.UUID)
	} else {
		log.Printf("[%s] no need login", client.UUID)
	}

	if result != "" || inRateID != "" {
		err = fmt.Errorf("calld(%s) already exists", callerId)
		return
	}

	var callerIP string

	if row.Via == "" && row.SrcIP != "" {
		callerIP = utils.ExtractIP(row.SrcIP)

	} else {
		callerIP = utils.ExtractIP(row.Via)
	}

	aniSip := sip.GetSipPart(row.ANI)
	dnisSip := sip.GetSipPart(row.DNIS)

	// 建立连接
	/*
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
	*/

	//log.Println("Login successfully!")

	//fmt.Printf("ani(%s) dnis(%s)\n", row[1], row[2])
	//command = fmt.Sprintf("call_simulation %s,5060,%s,%s", callerIP, aniSip, dnisSip)

	callerPort := "5060"

	// 构建命令
	command = fmt.Sprintf("call_simulation %s,%s,%s,%s\r\n", callerIP, callerPort, dnisSip, aniSip)
	log.Printf("[%s] Exec Command-> %s", callerId, command)

	content, err := client.CallSimulation(command)
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
	lrn := ""

	result = ""
	//
	if strings.Contains(content, "No Ingress Resource Found") {
		result = "No Ingress Resource Found"
		log.Printf("[%s]->result: %s", callerId, result)
	} else if strings.Contains(content, "Unauthorized IP Address") {
		result = "Unauthorized IP Address"
		log.Printf("[%s]->result: %s", callerId, result)
	} else if strings.Contains(content, "Ingress Rate Not Found") {
		result = "Ingress Rate Not Found"
		log.Printf("[%s]->result: %s", callerId, result)
	} else if strings.Contains(content, "YouMail Spam DB block") {
		result = "YouMail Spam DB block"
		log.Printf("[%s]->result: %s", callerId, result)
	} else {
		output, err := sip.ParseCallSimulationOutput(content)
		if err != nil {
			result = "Cannot Parse XMl"
			log.Printf("[%s]->result: %s", callerId, result)
		} else {
			if output.InBoundRate != "" {
				inbound_rate = output.InBoundRate
			}
			if output.InBoundRateId != "" {
				inbound_rate_id = output.InBoundRateId
			}
			if output.LRN != "" {
				lrn = output.LRN
			}
			//inbound_rate, inbound_rate_id, outbound_rate, outbound_rate_id, inTrunkId, outTrunkId = rate_utils.ParseRateFromContent(callerId, ani, dnis, aniSip, dnisSip, outVia, content)
			if !output.MatchRate(aniSip, dnisSip, outVia) { //未找到
				log.Printf("[call] CallerID(%s) ANI(%s) DNIS(%s) outVia(%s) not found out_bound", callerId, aniSip, dnisSip, outVia)
				result, err = convertor.ToJson(output)
				if err != nil {
					result = fmt.Sprintf("json err->%v", err)
				}
			} else { //找到
				inbound_rate = output.InBoundRate
				inbound_rate_id = output.InBoundRateId
				outbound_rate = output.OutBoundRate
				outbound_rate_id = output.OutBoundRateId
				inTrunkId = output.InTrunkId
				outTrunkId = output.OutTrunkId
				lrn = output.LRN

				log.Printf("[call] CallerID(%s) ANI(%s) DNIS(%s) LRN(%s) outVia(%s) inRate(%s) inRateId(%s) outRate(%s) outRateId(%s)", callerId, aniSip, dnisSip, lrn, outVia, inbound_rate, inbound_rate_id, outbound_rate, outbound_rate_id)

			}
		}

	}

	//_ = fmt.Sprintf("%s %s", inbound_rate, inbound_rate_id)
	ani := row.ANI
	dnis := row.DNIS
	row.ANI = sip.GetSipPart(ani)
	row.DNIS = sip.GetSipPart(dnis)

	row.Command = command
	row.Result = result
	row.LRN = lrn
	row.InRate = inbound_rate

	row.InRateID = inbound_rate_id
	row.OutRate = outbound_rate
	row.OutRateID = outbound_rate_id
	row.InTrunkId = inTrunkId
	row.OutTrunkId = outTrunkId

	defer pool.Put(client)

	return
}

func CalculateSipCost(path string) {
	connCount := 10
	// 创建连接池实例
	pool := telnet.NewTelnetClientPool(10)
	log.Printf("The telnet pool created with %d conns", connCount)

	csvFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	rows := []*csv_utils.PcapCsv{}

	if err := gocsv.UnmarshalFile(csvFile, &rows); err != nil { // Load clients from file
		panic(err)
	}

	n := 1
	all_count := len(rows)

	// 使用 WaitGroup 来等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 用 channel 控制最大并发数（限制为3个线程）
	sem := make(chan struct{}, 3) // 创建一个缓冲区大小为3的 channel

	for index, row := range rows {
		wg.Add(1)         // 增加等待计数
		sem <- struct{}{} // 获取一个信号量，限制并发数量

		go func(index int, pool *telnet.TelnetClientPool, row *csv_utils.PcapCsv) {
			defer wg.Done() // 完成时调用 Done

			err = handleRow(pool, row)
			if err != nil {
				log.Println("Skip row:", err)
				return
			}
			log.Printf("processing->%d/%d", n, all_count)

			rows[index] = row

			n += 1

			<-sem // 释放信号量
		}(index, pool, row) // 启动每个 goroutine

		//log.Printf("processing->%d/%d", n, all_count)

		//rows[index] = row

		/*
			fileName := filepath.Base(path)
			fileName = "res_" + fileName

			csvWriteFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}

			//每操作一次写入一次
			err = gocsv.MarshalFile(&rows, csvWriteFile) // Use this to save the CSV back to the file
			if err != nil {
				panic(err)
			}

			csvWriteFile.Close()

		*/
	}
}
