package xml_utils

import (
	"encoding/xml"
	"fmt"
	"log"
)

type OriginationLRNActionDNISAfter struct {
	XMLName xml.Name `xml:"Origination-LRN-Action-DNIS-After"`
	DNIS    string   `xml:"DNIS"`
	//LERG    string   `xml:"LERG"`
	//DNC     string   `xml:"DNC"`
}

func ParseOriginationLRNActionDNISAfter(data string) OriginationLRNActionDNISAfter {
	// 输入的 XML 数据
	// 创建结构体实例
	var rate OriginationLRNActionDNISAfter

	buff, err := sanitizeXML(data)
	if err != nil {
		fmt.Println(buff)
		log.Fatalf("[OriginationTrunk] XML sanitizeXML error: %v", err)
	}

	// 解析 XML 数据
	err = xml.Unmarshal([]byte(buff), &rate)
	if err != nil {
		fmt.Println(buff)
		log.Fatalf("[OriginationTrunk] XML Unmarshal error: %v", err)
	}

	// 打印解析结果
	//fmt.Printf("Parsed Struct: %+v\n", rate_utils)

	return rate
}
