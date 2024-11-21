package xml_utils

import (
	"encoding/xml"
	"log"
)

// 定义结构体与 XML 标签对应
type OriginationTrunkRate struct {
	Code          string `xml:"Code"`
	Country       string `xml:"Country"`
	CodeName      string `xml:"Code-Name"`
	RateID        string `xml:"Rate-ID"`
	Rate          string `xml:"Rate"`
	RateType      string `xml:"Rate-Type"`
	DnisType      string `xml:"Dnis-Type"`
	JurType       string `xml:"Jur-Type"`
	RateEffective string `xml:"Rate-Effective"`
}

func ParseOriginationTrunkRate(data string) OriginationTrunkRate {
	// 输入的 XML 数据
	// 创建结构体实例
	var rate OriginationTrunkRate

	// 解析 XML 数据
	err := xml.Unmarshal([]byte(data), &rate)
	if err != nil {
		log.Fatalf("XML Unmarshal error: %v", err)
	}

	// 打印解析结果
	//fmt.Printf("Parsed Struct: %+v\n", rate)

	return rate
}
