package sip

import (
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sip-parser/pkg/siprocket"
	"strings"
	"time"
)

func parsePacket(packet gopacket.Packet, timeStamp time.Time) (sipPacket *siprocket.SipMsg, err error) {
	err = errors.New("no sip msg")
	for _, layer := range packet.Layers() {
		switch layer.LayerType() {
		case layers.LayerTypeEthernet:
			//parseEthernetLayer(layer)
		case layers.LayerTypeIPv4:
			//parseIPV4Layer(layer)
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
			return input[start : start+end]
		}
	} else if strings.Contains(input, "sip:") {
		// 截取 <sip: 后的部分，去除后面的 @ 和 IP
		start := strings.Index(input, "sip:") + len("sip:")
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
