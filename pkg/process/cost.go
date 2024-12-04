package process

import (
	"encoding/csv"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gocarina/gocsv"
	"github.com/gogf/gf/v2/container/gtype"
	"io"
	"log"
	"os"
	"path/filepath"
	"sip-parser/pkg/sip"
	"sip-parser/pkg/utils"
	"sip-parser/pkg/utils/csv_utils"
	"sip-parser/pkg/utils/pool"
	"sip-parser/pkg/utils/telnet"
	"strings"
	"sync"
	"time"
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

func handleRow(pool pool.Pool, row *csv_utils.PcapCsv) (err error) {
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
	//inRateID := row.InRateID
	//inCost := row.InCost
	//outRate := row.OutRate
	//outRateID := row.OutRateID
	//outCost := row.OutCost
	inTrunkId := row.InTrunkId
	outTrunkId := row.OutTrunkId

	command := row.Command
	result := row.Result

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
	command = fmt.Sprintf("call_simulation %s,%s,%s,%s\r\n", callerIP, callerPort, aniSip, dnisSip)

	var conn interface{}
	var content string
	var client *telnet.TelnetClient
	for {
		//从连接池中取得一个连接
		conn, err = pool.Get()
		if conn == nil {
			continue
		}
		client = conn.(*telnet.TelnetClient)
		if err != nil {
			log.Println(err)
			continue
		}
		//client, err = pool.Get("127.0.0.1", "4320")

		if err != nil {
			log.Printf("put conn err->%v")
			pool.Close(client) //关闭进程
			continue
		}

		if !client.IsAuthentication {
			// 发送登录命令
			err = client.Login()
			if err != nil {
				log.Printf("login fail->%v", err)
				pool.Close(client) //关闭进程
				continue
			}
		}

		content, err = client.CallSimulation(command)
		if err != nil {
			err = fmt.Errorf("CallSimulation->%v", err)
			pool.Close(client) //关闭进程
			//client.Close() //关闭
			continue
		}
		pool.Close(client) //关闭进程
		if client.IsAvailable {
			defer pool.Put(client)
		}
		break
	}

	//_ = client.LoginOut()

	//
	inbound_rate := ""
	inbound_rate_id := ""
	outbound_rate := ""
	outbound_rate_id := ""
	lrn := ""

	result = ""
	isParseErr := false

	//
	if strings.Contains(content, "No Ingress Resource Found") {
		result = "No Ingress Resource Found"
		log.Printf("[%s]->command(%s) result: %s", callerId, command, result)
		isParseErr = true
	} else if strings.Contains(content, "Unauthorized IP Address") {
		result = "Unauthorized IP Address"
		log.Printf("[%s]->command(%s) result: %s", callerId, command, result)
		isParseErr = true
	} else if strings.Contains(content, "Ingress Rate Not Found") {
		result = "Ingress Rate Not Found"
		log.Printf("[%s]->command(%s) result: %s", callerId, command, result)
		isParseErr = true
	} else if strings.Contains(content, "YouMail Spam DB block") {
		result = "YouMail Spam DB block"
		log.Printf("[%s]->command(%s) result: %s", callerId, command, result)
		isParseErr = true
	} else if strings.Contains(content, "No Routing Plan Route") {
		result = "No Routing Plan Route"
		log.Printf("[%s]->command(%s) result: %s", callerId, command, result)
		isParseErr = true
	} else {
		log.Printf("[%s] Exec Command-> %s", callerId, command)

		output, err := sip.ParseCallSimulationOutput(content)
		if err != nil {
			result = "Cannot Parse XMl"
			log.Printf("[%s]->result: %s", callerId, result)
			isParseErr = true
			//return nil
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
			if output.InTrunkId != "" {
				inTrunkId = output.InTrunkId
			}
			//inbound_rate, inbound_rate_id, outbound_rate, outbound_rate_id, inTrunkId, outTrunkId = rate_utils.ParseRateFromContent(callerId, ani, dnis, aniSip, dnisSip, outVia, content)
			if !output.MatchRate(aniSip, dnisSip, outVia) { //未找到
				log.Printf("[call] CallerID(%s) ANI(%s) DNIS(%s) LRN(%s) InTrunkId(%s) outVia(%s) not found out_bound", callerId, aniSip, dnisSip, lrn, inTrunkId, outVia)
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

				log.Printf("[call] CallerID(%s) ANI(%s) DNIS(%s) LRN(%s) outVia(%s) inRate(%s) inRateId(%s) InTrunkId(%s) outRate(%s) outRateId(%s)", callerId, aniSip, dnisSip, lrn, outVia, inbound_rate, inbound_rate_id, inTrunkId, outbound_rate, outbound_rate_id)

			}
		}

	}

	if lrn == "" && inbound_rate == "" && !isParseErr {
		log.Printf("[call] CallerID(%s) cannnot get rate data, the content length is %d\n----\n%s", callerId, len(content), content)
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

	return
}

func CalculateSipCost(path string, costThreads int) {
	//connCount := 20
	// 创建连接池实例
	//pool := telnet.NewTelnetClientPool(connCount + 5)

	//factory 创建连接的方法
	factory := func() (interface{}, error) {
		//创建新的连接
		client := telnet.NewTelnetClient("127.0.0.1", "4320")
		err := client.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to create new telnet client: %v", err)
		}
		if !client.IsAuthentication {
			// 发送登录命令
			err = client.Login()
			if err != nil {
				return nil, fmt.Errorf("failed to login: %v", err)
			}
			//log.Printf("[%s] Successfully logged in!", client.UUID)
		} else {
			log.Printf("[%s] no need login", client.UUID)
		}

		return client, nil
	}
	//close 关闭连接的方法
	closeConn := func(v interface{}) error { return v.(*telnet.TelnetClient).CloseLogout() }

	//ping 检测连接的方法
	//ping := func(v interface{}) error { return v.(*telnet.TelnetClient).Ping() }

	log.Printf("The telnet pool created with %d conns", costThreads)
	initialCap := 5
	if initialCap >= costThreads {
		initialCap = costThreads
	}

	//创建一个连接池： 初始化5，最大空闲连接是20，最大并发连接30
	poolConfig := &pool.Config{
		InitialCap: initialCap,      //资源池初始连接数
		MaxIdle:    costThreads,     //最大空闲连接数
		MaxCap:     costThreads + 3, //最大并发连接数
		Factory:    factory,
		Close:      closeConn,
		//Ping:       ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 25 * time.Second,
	}

	connPool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		panic(err)
	}

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

	// 创建一个Int型的并发安全基本类型对象
	n := gtype.NewInt(1)

	all_count := len(rows)

	// 使用 WaitGroup 来等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 用 channel 控制最大并发数（限制为3个线程）
	sem := make(chan struct{}, costThreads) // 创建一个缓冲区大小为3的 channel

	is_need_write := false

	for index, row := range rows {
		wg.Add(1)         // 增加等待计数
		sem <- struct{}{} // 获取一个信号量，限制并发数量

		go func(index int, pool pool.Pool, row *csv_utils.PcapCsv) {
			defer wg.Done() // 完成时调用 Done

			defer n.Add(1)

			/*
				if strings.Contains(row.ANI, "%23") || strings.Contains(row.DNIS, "%23") {
					log.Printf("[%s] %s/%s contain # char", row.CallId, row.ANI, row.DNIS)
					//panic("11")
					err = handleRow(pool, row)
					if err != nil {
						log.Println("Skip row:", err)
						<-sem // 释放信号量
						return
					}
					log.Printf("processing->%d/%d", n.Val(), all_count)

					rows[index] = row
					//n.Add(1)
					is_need_write = true
					<-sem // 释放信号量
					return
				}

				if strings.Contains(row.ANI, "#23") || strings.Contains(row.DNIS, "#23") {
					row.ANI = strings.Replace(row.ANI, "#23", "#", -1)
					row.DNIS = strings.Replace(row.DNIS, "#23", "#", -1)

					log.Printf("%s/%s contain #23", row.ANI, row.DNIS)
					//panic("11")
					err = handleRow(pool, row)
					if err != nil {
						log.Println("Skip row:", err)
						<-sem // 释放信号量
						return
					}
					log.Printf("processing->%d/%d", n.Val(), all_count)

					rows[index] = row
					//n.Add(1)
					is_need_write = true
					<-sem // 释放信号量
					return
				}
			*/
			if row.InTrunkId != "" && row.InRate != "" && row.InRateID != "" {
				//InTrunkId不为空 不处理
				log.Printf("[%s] InTrunkId(%s) not empty, skip", row.CallId, row.InTrunkId)

			} else {
				if row.Result == "" {
					err = handleRow(pool, row)
					if err != nil {
						log.Println("Skip row:", err)
						<-sem // 释放信号量
						return
					}
					log.Printf("processing->%d/%d", n.Val(), all_count)

					rows[index] = row
					//n.Add(1)
					is_need_write = true
				} else {
					if row.InRate != "" {
						err = handleRow(pool, row)
						if err != nil {
							log.Println("Skip row:", err)
							<-sem // 释放信号量
							return
						}
						log.Printf("processing->%d/%d", n.Val(), all_count)

						rows[index] = row
						//n.Add(1)
						is_need_write = true
					} else {
						log.Printf("[%s] Result length(%d) has err, skip", row.CallId, len(row.Result))
					}
				}
			}

			<-sem // 释放信号量
		}(index, connPool, row) // 启动每个 goroutine

		if n.Val()%300 == 0 && is_need_write {
			log.Printf("saving data->%d/%d", n.Val(), all_count)
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

			is_need_write = false

			//panic("test")
		}

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
