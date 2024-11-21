package rate

import (
	"fmt"
	"log"
	"sip-parser/pkg/utils/xml_utils"
)

func ParseRateFromContent(callerID, content string) (inbound_rate, inbound_rate_id, outbound_rate, outbound_rate_id string) {
	// 将读取的内容转换为字符串
	//content := string(data)

	xmlList, err := xml_utils.ParseXMLToNodeList(content)
	if err != nil {
		log.Println("Error parsing XML:", err)
		return
	}
	//fmt.Println(xml_list, err)
	for _, nodes := range xmlList {
		// 检查 map 中是否包含指定的键
		originationTrunkKey := "Origination-Trunk-Rate"
		if data, exists := nodes[originationTrunkKey]; exists {
			raw_data := fmt.Sprintf("<Origination-Trunk-Rate>\n%s\n</Origination-Trunk-Rate>", data)
			//fmt.Printf("%s->%s", originationTrunkKey, raw_data)
			originationTrunk := xml_utils.ParseOriginationTrunkRate(raw_data)
			log.Printf("[inbound] CallerID(%s) RateID(%s) Rate(%s)\n", callerID, originationTrunk.RateID, originationTrunk.Rate)
			inbound_rate = originationTrunk.Rate
			inbound_rate_id = originationTrunk.RateID
		}

		terminationRouteKey := "Termination-Route"
		if data, exists := nodes[terminationRouteKey]; exists {
			raw_data := fmt.Sprintf("<Termination-Route>\n%s\n</Termination-Route>", data)
			//fmt.Printf("%s->%s", originationTrunkKey, raw_data)
			terminationRoute := xml_utils.ParseTrunkData(raw_data)
			//fmt.Println(terminationRoute.TerminationTrunk)
			if len(terminationRoute.TerminationTrunk) > 0 {
				for _, trunk := range terminationRoute.TerminationTrunk {
					FinalANI := trunk.FinalANI.ANI
					FinalANIReal := trunk.FinalANI.Real
					FinalDNIS := trunk.FinalDNIS.DNIS
					FinalDNISReal := trunk.FinalANI.Real

					log.Printf("[outbound] CallerID(%s) RateID(%s) Rate(%s) ANI(%s/%s) DNIS(%s/%s)\n", callerID, trunk.TrunkRate.RateID, trunk.TrunkRate.Rate, FinalANI, FinalANIReal, FinalDNIS, FinalDNISReal)

					outbound_rate = trunk.TrunkRate.Rate
					outbound_rate_id = trunk.TrunkRate.RateID
				}
			}
		}
	}

	return
}
