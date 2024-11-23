package sip

import (
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sip-parser/pkg/siprocket"
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
