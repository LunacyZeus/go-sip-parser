package sip

import (
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sip-parser/pkg/siprocket"
	"strings"
	"time"
)

func parsePacket(packet gopacket.Packet, timeStamp time.Time) (sipPacket *siprocket.SipMsg, srcIP string, destIP string, err error) {
	err = errors.New("no sip msg")
	for _, layer := range packet.Layers() {
		switch layer.LayerType() {
		case layers.LayerTypeEthernet:
			//parseEthernetLayer(layer)
		case layers.LayerTypeIPv4:
			//parseIPV4Layer(layer)
			// Parse IP Layer
			ip, _ := layer.(*layers.IPv4)
			srcIP = ip.SrcIP.String()
			destIP = ip.DstIP.String()
			//fmt.Println(srcIP, destIP)

		case layers.LayerTypeTCP:
			//parseTcpLayer(layer)
		case layers.LayerTypeUDP:
			err = nil
			udp, _ := layer.(*layers.UDP)

			sipPacket = siprocket.StreamParse(udp.Payload, timeStamp)
			//log.Println(string(sipPacket.Req.Src))
		default:
			//log.Println("- ", layer.LayerType())
		}
	}
	//parseEthernetLayer(packet)
	//parseIPV4Layer(packet)
	return
}

func ParsePart(part string) string {
	if strings.Contains(part, ";") {
		tmp := strings.Split(part, ";")
		return tmp[0]
	}
	if len(part) == 10 && part[0] != '1' {
		return "1" + part
	}
	return part
}

func GetSipPart(input string) string {
	if strings.Contains(input, "%") {
		input = strings.ReplaceAll(input, "%", "#")
	}
	// 去除 <sip: 和 @ 后面的内容
	if strings.Contains(input, "<sip:") {
		// 截取 <sip: 后的部分，去除后面的 @ 和 IP
		start := strings.Index(input, "<sip:") + len("<sip:")
		end := strings.Index(input[start:], "@")
		if end != -1 {
			// 打印号码部分
			return ParsePart(input[start : start+end])
		}
	} else if strings.Contains(input, "sip:") {
		// 截取 <sip: 后的部分，去除后面的 @ 和 IP
		start := strings.Index(input, "sip:") + len("sip:")
		end := strings.Index(input[start:], "@")
		if end != -1 {
			// 打印号码部分
			return ParsePart(input[start : start+end])
		}
	} else {
		// 如果不是 <sip: 格式，按 @ 分割并打印号码部分
		parts := strings.Split(input, "@")
		if len(parts) > 0 {
			return ParsePart(parts[0]) // 打印号码部分
		}
	}
	return ""
}
