package xml_utils

import (
	"encoding/xml"
	"fmt"
	"log"
)

// OriginationTrunk 用于映射 XML 中的 <Origination-Trunk> 元素
type OriginationTrunk struct {
	XMLName           xml.Name `xml:"Origination-Trunk"`
	CAP               string   `xml:"CAP"`
	CPS               string   `xml:"CPS"`
	CarrierName       string   `xml:"Carrier-Name"`
	TrunkID           string   `xml:"Trunk-ID"`
	TrunkName         string   `xml:"Trunk-Name"`
	RouteType         string   `xml:"Route-Type"`
	MediaType         string   `xml:"Media-Type"`
	ProfitMargin      string   `xml:"Profit-Margin"`
	ProfitType        string   `xml:"Profit-Type"`
	SHAKEN            SHAKEN   `xml:"SHAKEN"`
	StaticRouteName   string   `xml:"Static-Route-Name"`
	DynamicRouteName  string   `xml:"Dynamic-Route-Name"`
	RateTableName     string   `xml:"Rate-Table-Name"`
	RouteStrategyName string   `xml:"Route-Strategy-Name"`
}

func ParseOriginationTrunk(data string) OriginationTrunk {
	// 输入的 XML 数据
	// 创建结构体实例
	var rate OriginationTrunk

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
