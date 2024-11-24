package tests

import (
	"fmt"
	"sip-parser/pkg/process"
	"testing"
)

const inputFilepath = "data/test.pcap"

func TestGetCallSimulationCommand(t *testing.T) {
	rawLine := "<sip:+17049010476@162.212.247.96>;tag=sansay2659261434rdb103887"
	dnis := "<sip:+17049010476@162.212.247.96>;tag=sansay2659261434rdb103887"
	ani := "<sip:7356#+17049653278@172.241.26.23>"

	callerIP := process.ExtractIP(rawLine)
	callerPort := "5060"
	aniSip := process.GetSipPart(dnis)
	dnisSip := process.GetSipPart(ani)

	// 构建命令
	command := fmt.Sprintf("call_simulation %s,%s,%s,%s\r\n", callerIP, callerPort, aniSip, dnisSip)

	t.Logf("%s", command)
}
