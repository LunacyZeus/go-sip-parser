package process

import (
	"fmt"
	"os"
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

func StartTelnet(csvFilePath string) {
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
	callerIp := "207.223.71.199"
	callerPort := "5060"
	ani := "+13134638035"
	dnis := "+13133854865"

	result, err := client.CallSimulation(callerIp, callerPort, ani, dnis)
	if err != nil {
		fmt.Println("CallSimulation", err)
		return
	}

	writeToFile("call.xml", result)
}
