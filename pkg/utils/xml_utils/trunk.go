package xml_utils

import (
	"encoding/xml"
	"fmt"
	"log"
)

// Define the struct types to match the XML structure
type TerminationRoute struct {
	XMLName          xml.Name           `xml:"Termination-Route"`
	StaticRouteName  string             `xml:"Static-Route-Name"`
	StaticPrefix     string             `xml:"Static-Prefix"`
	RouteStrategy    string             `xml:"Route-Strategy"`
	TerminationTrunk []TerminationTrunk `xml:"Termination-Trunk"`
}

type TerminationTrunk struct {
	CarrierName          string            `xml:"Carrier-Name"`
	TrunkID              string            `xml:"Trunk-ID"`
	TrunkName            string            `xml:"Trunk-Name"`
	RateTableName        string            `xml:"Rate-Table-Name"`
	CPS                  string            `xml:"CPS"`
	CAP                  string            `xml:"CAP"`
	BillAfterAction      string            `xml:"Bill-After-Action"`
	Strategy             string            `xml:"Strategy"`
	SHAKEN               SHAKEN            `xml:"SHAKEN"`
	FinalANI             FinalANI          `xml:"Final-ANI"`
	FinalDNIS            FinalDNIS         `xml:"Final-DNIS"`
	ASR                  string            `xml:"ASR"`
	ACD                  string            `xml:"ACD"`
	TrunkRate            TrunkRate         `xml:"Trunk-Rate"`
	TerminationHost      []TerminationHost `xml:"Termination-Host"`
	TerminationSignature string            `xml:"Termination-Signature"`
	PHeaders             string            `xml:"P-Headers"`
}

type SHAKEN struct {
	VfyPolicy string `xml:"Vfy-policy"`
}

type FinalANI struct {
	Actions []Action `xml:"Actions"`
	ANI     string   `xml:"ANI"`
	Real    string   `xml:"Real"`
	DNO     []string `xml:"DNO"`
}

type FinalDNIS struct {
	Actions []Action `xml:"Actions"`
	DNIS    string   `xml:"DNIS"`
	Real    string   `xml:"Real"`
	DNC     string   `xml:"DNC"`
}

type Action struct {
	Action string `xml:"Action"`
	Digit  string `xml:"Digit"`
}

type TrunkRate struct {
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

type TerminationHost struct {
	CPS    string `xml:"CPS"`
	CAP    string `xml:"CAP"`
	Port   string `xml:"Port"`
	HostIP string `xml:"Host-IP"`
}

func ParseTrunkData(data string) TerminationRoute {
	buff, err := sanitizeXML(data)
	if err != nil {
		fmt.Println(data)
		log.Fatalf("[TerminationRoute] XML sanitizeXML error: %v", err)
	}

	var route TerminationRoute
	err = xml.Unmarshal([]byte(buff), &route)
	if err != nil {
		log.Fatalf("[TerminationRoute] Error unmarshalling XML: %v", err)
	}

	//fmt.Printf("Static Route Name: %s\n", route.StaticRouteName)
	//fmt.Printf("Static Prefix: %s\n", route.StaticPrefix)
	//fmt.Printf("Route Strategy: %s\n", route.RouteStrategy)
	//fmt.Printf("Carrier Name: %s\n", route.TerminationTrunk.CarrierName)
	// Add more print statements as needed to access other fields

	return route
}
