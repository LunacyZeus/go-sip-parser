package rate_utils

import (
	"fmt"
	"log"
	"sip-parser/pkg/utils"
	"sip-parser/pkg/utils/xml_utils"
	"strings"
)

func removePlusPrefix(s string) string {
	if strings.HasPrefix(s, "+") {
		return s[1:] // 去掉第一个字符
	}
	return s // 原样返回
}

func ParseRateFromContent(callerID, ani, dnis, aniSip, dnisSip, outVia, content string) (inbound_rate, inbound_rate_id, outbound_rate, outbound_rate_id string) {
	// 将读取的内容转换为字符串
	//content := string(data)

	aniSip = removePlusPrefix(aniSip)
	dnisSip = removePlusPrefix(dnisSip)

	xmlList, err := xml_utils.ParseXMLToNodeList(content)
	if err != nil {
		log.Println("Error parsing XML:", err)
		return
	}
	//fmt.Println(xml_list, err)
	notMatchMsgs := []string{}
	isFound := false

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
					hosts := []string{}
					for _, host := range trunk.TerminationHost {
						hosts = append(hosts, host.HostIP)
					}

					hostsStr := strings.Join(hosts, ",")

					FinalANI := trunk.FinalANI.ANI
					//FinalANIReal := trunk.FinalANI.Real
					FinalDNIS := trunk.FinalDNIS.DNIS
					//FinalDNISReal := trunk.FinalANI.Real

					if dnisSip == removePlusPrefix(FinalANI) && (aniSip == removePlusPrefix(FinalDNIS)) {
						outbound_rate = trunk.TrunkRate.Rate
						outbound_rate_id = trunk.TrunkRate.RateID
						log.Printf("[outbound] CallerID(%s) RateID(%s) Rate(%s) ANI(%s) DNIS(%s) host(%s) outVia(%s)\n", callerID, trunk.TrunkRate.RateID, trunk.TrunkRate.Rate, FinalANI, FinalDNIS, hostsStr, outVia)
						isFound = true
						return
					} else {
						viaIP := utils.ExtractIP(outVia)
						if strings.Contains(hostsStr, viaIP) {
							outbound_rate = trunk.TrunkRate.Rate
							outbound_rate_id = trunk.TrunkRate.RateID
							log.Printf("[outbound] CallerID(%s) RateID(%s) Rate(%s) ANI(%s) DNIS(%s) host(%s) outVia(%s)\n", callerID, trunk.TrunkRate.RateID, trunk.TrunkRate.Rate, FinalANI, FinalDNIS, hostsStr, outVia)
							isFound = true
							return
						}
						msg := fmt.Sprintf("[outbound] CallerID(%s) RateID(%s) Rate(%s) ANI(%s!=%s) DNIS(%s!=%s) host(%s) outVia(%s) not match\n", callerID, trunk.TrunkRate.RateID, trunk.TrunkRate.Rate, dnisSip, removePlusPrefix(FinalANI), aniSip, removePlusPrefix(FinalDNIS), hostsStr, outVia)
						notMatchMsgs = append(notMatchMsgs, msg)
					}
				}
			}
		}
	}

	if !isFound { //没搜索到
		pMsg := strings.Join(notMatchMsgs, "\n")
		log.Printf("%s", pMsg)
	}

	return
}
