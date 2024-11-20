package process

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func sendCommand(conn net.Conn, command string, reader *bufio.Reader) error {
	// 发送命令
	fmt.Printf("Sending command: %s", command)
	_, err := conn.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	// 读取服务器响应
	fmt.Println("Waiting for response...")
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // 设置超时时间
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	fmt.Printf("Server response: %s", strings.TrimSpace(response))
	return nil
}

func StartTelnet(csvFilePath string) {
	log.Println(csvFilePath)

	// 连接到 127.0.0.1:4320
	conn, err := net.Dial("tcp", "127.0.0.1:4320")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to 127.0.0.1:4320")

	// 创建读取器，用于接收服务器返回的数据
	reader := bufio.NewReader(conn)

	// 发送 login 命令
	if err := sendCommand(conn, "login\n", reader); err != nil {
		fmt.Println("Error during login:", err)
		os.Exit(1)
	}

	// 等待服务器响应后发送 call_simulation
	if err := sendCommand(conn, "call_simulation\n", reader); err != nil {
		fmt.Println("Error during call_simulation:", err)
		os.Exit(1)
	}
}
