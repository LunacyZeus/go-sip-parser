package tests

import (
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/fileutil"
	"log"
	"sip-parser/pkg/sip"
	"testing"
)

const outputFilepath = "../data/call.xml"

func TestParseCallSimulationOutput(t *testing.T) {
	// 读取 XML 文件
	xmlFile, err := fileutil.ReadFileToString(outputFilepath)
	if err != nil {
		log.Fatalf("无法读取文件: %v", err)
	}

	output, err := sip.ParseCallSimulationOutput(xmlFile)
	if err != nil {
		t.Fatalf("无法解析 XML 文件: %v", err)
	}
	is_ok := output.MatchRate("131338548651", "13134638035", "sip:207.38.67.251")

	if !is_ok {
		result, err := convertor.ToJson(output)

		if err != nil {
			fmt.Printf("%v", err)
		}
		t.Log(result)

	}

	t.Logf("output->%v", output)
}
