package process

import (
	"fmt"
	"log"
	"os"
	"sip-parser/pkg/sip"
)

func FilterPcapFileOrFolder(path, callId string) {
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

		fp, err := sip.LoadSIPTraceFromPcap(path)
		if err != nil {
			log.Panic(err)
		}

		// Search the SIP packets for the filters
		sip.HandleSipPackets(manager, fp)
	}

	log.Printf("筛选callId->%s", callId)
	session, exists := manager.GetSession(callId)
	if !exists {
		log.Fatal("获取session失败")
	}
	for _, msg := range session.Messages {
		fmt.Println(msg.String())
	}
	log.Printf("%d msgs got", len(session.Messages))
}
