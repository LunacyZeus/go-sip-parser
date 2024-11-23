package process

import (
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"sip-parser/pkg/siprocket"
	"time"
)

var (
	pcapFile string = "data/test/202411140021.pcap"
	handle   *pcap.Handle
	err      error
)

func parsePacket(packet gopacket.Packet) (sipPacket siprocket.SipMsg, err error) {
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
			sipPacket = siprocket.Parse(udp.Payload, time.Now())
			//log.Println(string(sipPacket.To.Src))
		default:
			//log.Println("- ", layer.LayerType())
		}
	}
	//parseEthernetLayer(packet)
	//parseIPV4Layer(packet)
	return
}

func parseEthernetLayer(layer gopacket.Layer) {
	log.Println("Ethernet layer detected.")
	ethernetPacket, _ := layer.(*layers.Ethernet)
	log.Println("Source MAC: ", ethernetPacket.SrcMAC)
	log.Println("Destination MAC: ", ethernetPacket.DstMAC)
	log.Println("Ethernet type: ", ethernetPacket.EthernetType)

}

func parseIPV4Layer(layer gopacket.Layer) {
	log.Println("IPv4 layer detected.")
	ip, _ := layer.(*layers.IPv4)

	log.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
	log.Println("Protocol: ", ip.Protocol)
}

func parseTcpLayer(layer gopacket.Layer) {
	log.Println("TCP layer detected.")
	tcp, _ := layer.(*layers.TCP)

	log.Printf("From %s to %s\n", tcp.SrcPort, tcp.DstPort)

	//tcpassembly.NewAssembler(nil)

}

func parseUdpLayer(layer gopacket.Layer) {
	//log.Println("UDP layer detected.")
	udp, _ := layer.(*layers.UDP)

	//log.Printf("From %s to %s\n", udp.SrcPort, udp.DstPort)
	//log.Printf("UDP payload is %s \n", string(udp.Payload))

	sipPacket := siprocket.Parse(udp.Payload, time.Now())
	log.Println(string(sipPacket.To.Src))

}

func TestFunc() {
	// 打开PCAP文件，而不是设备
	handle, err = pcap.OpenOffline(pcapFile)
	// 如果打开文件时发生错误，记录错误并终止程序
	if err != nil {
		log.Fatal(err)
	}
	// 确保在函数结束时关闭文件句柄
	defer handle.Close()

	// 创建一个新的包源，用于从文件中读取包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// 遍历文件中的所有包
	for packet := range packetSource.Packets() {
		// 处理每个包，这里简单地打印出来
		parsePacket(packet)
	}
}
