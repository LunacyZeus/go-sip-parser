package sip

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marv2097/siprocket"
	"sip-parser/pkg/gopcap"
)

type SipMessage struct {
	pct       siprocket.SipMsg
	Timestamp time.Duration
}

func ParsePcapFile(file string) (gopcap.PcapFile, error) {
	if file == "" {
		return gopcap.PcapFile{}, errors.New("empty file specified")
	}

	pcapfile, _ := os.Open(file)
	parsed, err := gopcap.Parse(pcapfile)
	if err != nil {
		return gopcap.PcapFile{}, fmt.Errorf("cannot parse the pcap file: %s", err)
	}
	return parsed, nil
}

func ParseSIPTrace(trace gopcap.PcapFile) ([]SipMessage, error) {
	var results []SipMessage
	for _, packet := range trace.Packets {
		var r SipMessage
		d := packet.Data
		if d == nil {
			continue
		}

		td := d.LinkData().InternetData().TransportData()
		if td == nil {
			log.Println("unexpected transport data")
			continue
		}

		sipPacket := siprocket.Parse(td)
		r.pct = sipPacket
		r.Timestamp = packet.Timestamp
		results = append(results, r)
	}
	return results, nil
}

func HandleSipPackets(sipPackets []SipMessage) {
	// 创建一个 SIP 会话管理器
	manager := NewSipSessionManager()

	for _, sipp := range sipPackets {
		//fmt.Println(sipp.timestamp.Microseconds(), sipp.pct)
		//fmt.Println(sipp.Timestamp.Microseconds(), string(sipp.pct.CallId.Src), string(sipp.pct.Cseq.Src), string(sipp.pct.From.Src), string(sipp.pct.To.Src))
		//siprocket.PrintSipStruct(&sipp.pct)

		callID := string(sipp.pct.CallId.Value)

		// 如果 Call-ID 存在，则处理该会话
		if callID != "" {
			//fmt.Println(callID)
			// 查找会话，如果不存在则创建新会话
			session, exists := manager.GetSession(callID)
			if !exists {
				// 创建新的 SIP 会话
				session = NewSipSession(callID)
				manager.AddSession(session)
			}

			// 添加消息到会话
			session.AddMessage(&sipp)
		}
	}

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
		for _, session := range sessions {
			if session.Status == COMPLETED { //只解析成功的
				fmt.Println(session.String())
			}
		}
	}

	return
}

func HandleSipPackets1(sipPackets []SipMessage) {
	for _, sipp := range sipPackets {
		//fmt.Println(sipp.timestamp.Microseconds(), sipp.pct)
		//fmt.Println(sipp.Timestamp.Microseconds(), string(sipp.pct.CallId.Src), string(sipp.pct.Cseq.Src), string(sipp.pct.From.Src), string(sipp.pct.To.Src))
		//siprocket.PrintSipStruct(&sipp.pct)

		callID := string(sipp.pct.CallId.Value)

		// 如果 Call-ID 存在，则处理该会话
		if callID != "" {
			fmt.Println(callID)
		}
	}

	return
}
