package sip

import (
	"fmt"
	"sip-parser/pkg/utils"
	"sip-parser/pkg/utils/xml_utils"
	"strings"
)

type TerminationRoute struct {
	Hosts         string
	FinalANI      string
	FinalANIReal  string
	FinalDNIS     string
	FinalDNISReal string
	Rate          string
	RateId        string
	TrunkId       string
}

type SimOutput struct {
	InTrunkId      string
	OutTrunkId     string
	InBoundRate    string
	InBoundRateId  string
	OutBoundRate   string
	OutBoundRateId string
	LRN            string
	Routes         []TerminationRoute
}

func removePlusPrefix(s string) string {
	if strings.HasPrefix(s, "+") {
		return s[1:] // 去掉第一个字符
	}
	return s // 原样返回
}

func (o *SimOutput) MatchRate(aniSip, dnisSip, outVia string) bool {
	for _, route := range o.Routes {
		viaIP := utils.ExtractIP(outVia)

		if strings.Contains(dnisSip, removePlusPrefix(route.FinalANI)) && strings.Contains(aniSip, removePlusPrefix(route.FinalDNIS)) && strings.Contains(route.Hosts, viaIP) {
			o.OutBoundRate = route.Rate
			o.OutBoundRateId = route.RateId
			o.OutTrunkId = route.TrunkId
			return true
		} else {
			if strings.Contains(route.Hosts, viaIP) {
				o.OutBoundRate = route.Rate
				o.OutBoundRateId = route.RateId
				o.OutTrunkId = route.TrunkId
				return true
			}
			//msg := fmt.Sprintf("[outbound] CallerID(%s) RateID(%s) Rate(%s) ANI(%s!=%s) DNIS(%s!=%s) host(%s) outVia(%s) not match\n", callerID, trunk.TrunkRate.RateID, trunk.TrunkRate.Rate, dnisSip, removePlusPrefix(FinalANI), aniSip, removePlusPrefix(FinalDNIS), hostsStr, outVia)
			//notMatchMsgs = append(notMatchMsgs, msg)
		}
	}
	return false
}

func ParseCallSimulationOutput(xmlContent string) (output *SimOutput, err error) {
	output = new(SimOutput)
	xmlList, err := xml_utils.ParseXMLToNodeList(xmlContent)
	if err != nil {
		//log.Println("Error parsing XML:", err)
		return
	}

	for _, nodes := range xmlList {
		// 检查 map 中是否包含指定的键
		originationTrunkKey := "Origination-Trunk"
		if data, exists := nodes[originationTrunkKey]; exists {
			raw_data := fmt.Sprintf("<Origination-Trunk>\n%s\n</Origination-Trunk>", data)
			//fmt.Printf("%s->%s", originationTrunkKey, raw_data)
			originationTrunk := xml_utils.ParseOriginationTrunk(raw_data)
			output.InTrunkId = originationTrunk.TrunkID
		}

		originationLRNActionAfterKey := "Origination-LRN-Action-DNIS-After"
		if data, exists := nodes[originationLRNActionAfterKey]; exists {
			raw_data := fmt.Sprintf("<Origination-LRN-Action-DNIS-After>\n%s\n</Origination-LRN-Action-DNIS-After>", data)
			//fmt.Printf("%s->%s", originationTrunkKey, raw_data)
			originationLRNActionAfter := xml_utils.ParseOriginationLRNActionDNISAfter(raw_data)
			output.LRN = originationLRNActionAfter.DNIS
		}

		originationTrunkRateKey := "Origination-Trunk-Rate"
		if data, exists := nodes[originationTrunkRateKey]; exists {
			raw_data := fmt.Sprintf("<Origination-Trunk-Rate>\n%s\n</Origination-Trunk-Rate>", data)
			//fmt.Printf("%s->%s", originationTrunkKey, raw_data)
			originationTrunkRate := xml_utils.ParseOriginationTrunkRate(raw_data)
			//log.Printf("[inbound] CallerID(%s) RateID(%s) Rate(%s)\n", callerID, originationTrunkRate.RateID, originationTrunkRate.Rate)
			output.InBoundRate = originationTrunkRate.Rate
			output.InBoundRateId = originationTrunkRate.RateID
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

					route := TerminationRoute{
						Hosts:         hostsStr,
						FinalANI:      trunk.FinalANI.ANI,
						FinalANIReal:  trunk.FinalANI.Real,
						FinalDNIS:     trunk.FinalDNIS.DNIS,
						FinalDNISReal: trunk.FinalDNIS.Real,
						Rate:          trunk.TrunkRate.Rate,
						RateId:        trunk.TrunkRate.RateID,
						TrunkId:       trunk.TrunkID,
					}

					output.Routes = append(output.Routes, route)
				}
			}
		}
	}

	return
}
