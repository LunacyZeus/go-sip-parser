package process

import (
	"fmt"
	"os"
	"sip-parser/pkg/utils/rate"
	"sip-parser/pkg/utils/telnet"
)

// 将字符串写入文件
func writeToFile(filename, content string) error {
	file, err := os.Create(filename) // 创建文件，若已存在会覆盖
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content) // 写入内容
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}

	fmt.Printf("Successfully wrote to file: %s\n", filename)
	return nil
}

type CallSimulationParams struct {
	CallerIp   string
	CallerPort string
	Ani        string
	Dnis       string
}

func StartTelnet(params CallSimulationParams) {
	// 创建客户端实例
	client := telnet.NewTelnetClient("127.0.0.1", "4320")

	// 建立连接
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	// 发送登录命令
	err = client.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Login successfully!")
	callerIp := params.CallerIp
	callerPort := params.CallerPort
	ani := params.Ani
	dnis := params.Dnis

	content, err := client.CallSimulation(callerIp, callerPort, ani, dnis)
	if err != nil {
		fmt.Println("CallSimulation", err)
		return
	}

	//No Ingress Resource Found

	writeToFile("call1.xml", content)

	rate.ParseRateFromContent("", content)

}
