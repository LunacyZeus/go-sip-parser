package process

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sip-parser/pkg/sip"
	"sip-parser/pkg/utils/csv_utils"
	"strings"
)

var manager *sip.SipSessionManager

func ProcessFileOrFolder(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error accessing path: %v\n", err)
		return
	}

	// 创建一个 SIP 会话管理器
	manager = sip.NewSipSessionManager()

	if info.IsDir() {
		fmt.Printf("Processing folder: %s\n", path)
		processFolder(path)
	} else {
		fmt.Printf("Processing file: %s\n", path)

		fileName := filepath.Base(path)

		fp, err := sip.LoadSIPTraceFromPcap(path)
		if err != nil {
			log.Panic(err)
		}

		// Search the SIP packets for the filters
		sip.HandleSipPackets(manager, fp)

		sessions := manager.Sessions

		saveCsvFileName := fmt.Sprintf("%s.csv", fileName)
		// 写入数据
		csv_utils.SaveDataCsv(saveCsvFileName, sessions)
	}
}

func processFolder(folderPath string) {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".pcap") {
			fmt.Printf("Found pcap file: %s\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error processing folder: %v\n", err)
	}
}