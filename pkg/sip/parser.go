package sip

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"time"

	"sip-parser/pkg/gopcap"
	"sip-parser/pkg/siprocket"
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

func LoadSIPTraceFromPcap(file string) ([]*siprocket.SipMsg, error) {
	if file == "" {
		return nil, errors.New("empty file specified")
	}

	pcapfile, _ := os.Open(file)
	trace, err := gopcap.Parse(pcapfile)
	if err != nil {
		return nil, fmt.Errorf("cannot parse the pcap file: %s", err)
	}
	pcapfile.Close()

	var results []*siprocket.SipMsg
	for _, packet := range trace.Packets {
		d := packet.Data
		if d == nil {
			continue
		}

		td := d.LinkData().InternetData().TransportData()
		if td == nil {
			log.Println("unexpected transport data")
			continue
		}

		sipPacket := siprocket.Parse(td, time.Now())

		results = append(results, &sipPacket)
	}
	return results, nil
}

func ParseSIPTrace(trace gopcap.PcapFile) ([]*siprocket.SipMsg, error) {
	var results []*siprocket.SipMsg
	for _, packet := range trace.Packets {
		d := packet.Data
		if d == nil {
			continue
		}

		td := d.LinkData().InternetData().TransportData()
		if td == nil {
			log.Println("unexpected transport data")
			continue
		}

		sipPacket := siprocket.Parse(td, time.Now())

		results = append(results, &sipPacket)
	}
	return results, nil
}

func HandleSipPackets(manager *SipSessionManager, sipPackets []*siprocket.SipMsg) {
	for _, sipp := range sipPackets {
		//fmt.Println(sipp.timestamp.Microseconds(), sipp.pct)
		//fmt.Println(sipp.Timestamp.Microseconds(), string(sipp.pct.CallId.Src), string(sipp.pct.Cseq.Src), string(sipp.pct.From.Src), string(sipp.pct.To.Src))
		//siprocket.PrintSipStruct(&sipp.pct)

		callID := string(sipp.CallId.Value)

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
			session.AddMessage(sipp)
			manager.LatestPktTimestamp = sipp.Timestamp
		}
	}

	return
}

func HandleSipPackets1(sipPackets []*siprocket.SipMsg) {
	for _, sipp := range sipPackets {
		//fmt.Println(sipp.timestamp.Microseconds(), sipp.pct)
		//fmt.Println(sipp.Timestamp.Microseconds(), string(sipp.pct.CallId.Src), string(sipp.pct.Cseq.Src), string(sipp.pct.From.Src), string(sipp.pct.To.Src))
		//siprocket.PrintSipStruct(&sipp.pct)

		callID := string(sipp.CallId.Value)

		// 如果 Call-ID 存在，则处理该会话
		if callID != "" {
			//fmt.Println(callID)
		}
	}
	fmt.Println("done")

	return
}

func LoadSIPTraceFromPcapStream(file string) ([]*siprocket.SipMsg, error) {
	// 打开PCAP文件，而不是设备
	handle, err := pcap.OpenOffline(file)
	// 如果打开文件时发生错误，记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}
	// 确保在函数结束时关闭文件句柄
	defer handle.Close()

	var results []*siprocket.SipMsg

	// 创建一个新的包源，用于从文件中读取包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// 遍历文件中的所有包
	for packet := range packetSource.Packets() {
		// 处理每个包，这里简单地打印出来
		pkt, srcIP, destIP, err := parsePacket(packet, packet.Metadata().Timestamp)
		pkt.SrcIP = srcIP
		pkt.DestIP = destIP

		if err != nil {
			continue
		} else {
			results = append(results, pkt)
		}
	}
	return results, nil
}
