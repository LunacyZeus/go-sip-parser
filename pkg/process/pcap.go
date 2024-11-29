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

func LoadPcap(path string) {
	ProcessFileOrFolder(path)
}

func ProcessFileOrFolder(path string) {
	log.Printf("Parsing pcap->%s", path)
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

		//folderPathFileName := strings.ReplaceAll(path, "/", "-")
		//sessions := manager.Sessions
		//saveCsvFileName := fmt.Sprintf("%s.csv", folderPathFileName)
		// 写入数据
		//csv_utils.SaveDataCsv(saveCsvFileName, sessions)

	} else {
		fmt.Printf("Processing file: %s\n", path)

		fp, err := sip.LoadSIPTraceFromPcapStream(path)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("pcap file(%s) loaded", path)

		// Search the SIP packets for the filters
		sip.HandleSipPackets(manager, fp)
	}
}

func processFolder(folderPath string) {
	all_count := 0
	n := 0
	fileName := filepath.Base(folderPath)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".pcap") {
			log.Printf("Found and parsing pcap file: %s\n", path)
			//fileName := filepath.Base(path)

			fp, err := sip.LoadSIPTraceFromPcapStream(path)
			if err != nil {
				log.Printf("cannot parsing file: %s err:%v", path, err)
				return nil
			}
			log.Printf("pcap file(%s) loaded", path)

			// Search the SIP packets for the filters
			sip.HandleSipPackets(manager, fp)

			manager.Statistics()
			manager.MatchCall()
			sessions := manager.GetAndDeleteAllCompleteCall(manager.LatestPktTimestamp.UnixMilli())

			saveCsvFileName := fmt.Sprintf("%s.csv", fileName)

			all_count += len(sessions)

			if n == 1 { //第一次 写入文件
				// 全新写入数据
				csv_utils.SaveDataCsv(saveCsvFileName, sessions)
			} else {
				// 追加写入数据
				csv_utils.AppendDataCsv(saveCsvFileName, sessions)
			}

			n += 1
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error processing folder: %v\n", err)
		return
	}

	fmt.Printf("%d CSV data wirte successfully\n", all_count)
}
