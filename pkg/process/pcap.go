package process

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sip-parser/pkg/sip"
	"sip-parser/pkg/utils/csv_utils"
	"strings"
	"time"
)

var manager *sip.SipSessionManager

func LoadPcap(path string) {
	//f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	//defer f.Close()
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	startT := time.Now() //计算当前时间
	ProcessFileOrFolder(path)
	tc := time.Since(startT) //计算耗时
	fmt.Printf("time cost = %v\n", tc)

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

		log.Printf("pcap file(%s) get csv to restore", path)
		sessions := manager.GetAndDeleteAllCompleteCall(manager.LatestPktTimestamp.UnixMilli())

		saveCsvFileName := fmt.Sprintf("%s.csv", path)

		all_count := len(sessions)

		if strings.Contains(saveCsvFileName, "/") {
			saveCsvFileName = strings.ReplaceAll(saveCsvFileName, "/", "_")
		}
		// 全新写入数据
		csv_utils.SaveDataCsv(saveCsvFileName, sessions)
		fmt.Printf("%s %d CSV data wirte successfully\n", saveCsvFileName, all_count)
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

			log.Printf("pcap file(%s) handle pcap", path)
			err := sip.LoadSIPTraceFromPcapStreamWithManager(manager, path)
			if err != nil {
				log.Printf("cannot parsing file: %s err:%v", path, err)
				return nil
			}
			log.Printf("pcap file(%s) loaded", path)

			//fp, err := sip.LoadSIPTraceFromPcapStream(path)
			//if err != nil {
			//	log.Panic(err)
			//}
			//log.Printf("pcap file(%s) loaded", path)

			// Search the SIP packets for the filters
			//sip.HandleSipPackets(manager, fp)

			// Search the SIP packets for the filters
			//log.Printf("pcap file(%s) handle packets", path)
			//sip.HandleSipPackets(manager, fp)

			log.Printf("pcap file(%s) Statistics", path)
			manager.Statistics()
			log.Printf("pcap file(%s) MatchCall", path)

			manager.MatchCall()
			log.Printf("pcap file(%s) get csv to restore", path)
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

	fmt.Printf("%s %d CSV data wirte successfully\n", fileName, all_count)
}
