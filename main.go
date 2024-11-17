package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sip-parser/pkg/sip"
	"sip-parser/pkg/utils/csv_utils"
	"strings"
)

var (
	file   = flag.String("file", "", "The SIP pcap file that will be parsed")
	to     = flag.String("to", "", "SIP To: Header")
	from   = flag.String("from", "", "SIP From: Header")
	callid = flag.String("callid", "", "SIP Call-ID header")
	help   = flag.Bool("help", false, "Display usage help")
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:8081", nil))
	}()

	flag.Usage = func() {
		fmt.Printf("Usage: %s \n", os.Args[0])
		fmt.Println("If no `--to` and `--from` are specified then the program will output `To:` and `From:` from all SIP dialogs.")
		fmt.Println("Parameters:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *file == "" {
		errorOut("you need to specify the file with --file")
	}

	*from = strings.ToLower(*from)
	*to = strings.ToLower(*to)
	*callid = strings.ToLower(*callid)

	fp, err := sip.LoadSIPTraceFromPcap(*file)
	if err != nil {
		errorOut(err.Error())
	}

	// 创建一个 SIP 会话管理器
	manager := sip.NewSipSessionManager()

	// Search the SIP packets for the filters
	sip.HandleSipPackets(manager, fp)

	callId := ""
	if callId != "" {
		session, exists := manager.GetSession(callId)
		if !exists {
			log.Fatal("获取session失败")
		}
		for _, msg := range session.Messages {
			fmt.Println(msg.String())
		}
		fmt.Println(session.String())
	} else {
		sessions := manager.Sessions

		// 写入数据
		csv_utils.SaveDataCsv("data.csv", sessions)
		/*
			for _, session := range sessions {
				if session.Status == sip.COMPLETED { //只解析成功的
					fmt.Println(session.String())
				}
			}

		*/
	}

	select {}
}
