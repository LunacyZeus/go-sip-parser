package xml_utils

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
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

func sanitizeXML(input string) (string, error) {
	// 将未转义的 & 替换为 &amp;
	return strings.ReplaceAll(input, "&", "&amp;"), nil
}

func ParseOriginationTrunkRate(data string) OriginationTrunkRate {
	// 输入的 XML 数据
	// 创建结构体实例
	var rate OriginationTrunkRate

	buff, err := sanitizeXML(data)
	if err != nil {
		fmt.Println(buff)
		log.Fatalf("[OriginationTrunkRate] XML sanitizeXML error: %v", err)
	}

	// 解析 XML 数据
	err = xml.Unmarshal([]byte(buff), &rate)
	if err != nil {
		fmt.Println(buff)
		log.Fatalf("[OriginationTrunkRate] XML Unmarshal error: %v", err)
	}

	// 打印解析结果
	//fmt.Printf("Parsed Struct: %+v\n", rate)

	return rate
}
