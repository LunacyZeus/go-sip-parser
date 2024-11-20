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

func sendLoginCommand(conn net.Conn, reader *bufio.Reader) (bool, error) {
	command := "login\r\n"
	// 发送命令
	fmt.Printf("Sending: %s", command)
	_, err := conn.Write([]byte(command))
	if err != nil {
		return false, fmt.Errorf("failed to send command: %w", err)
	}

	// 读取服务器响应
	fmt.Println("Waiting for Login response...")
	conn.SetReadDeadline(time.Now().Add(8 * time.Second)) // 设置超时时间
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}
	fmt.Printf("login recv: %s", strings.TrimSpace(response))
	return true, nil
}

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
	_, err = sendLoginCommand(conn, reader)
	if err != nil {
		fmt.Println("Error during login:", err)
		os.Exit(1)
	}

	// 等待服务器响应后发送 call_simulation
	call_simulation_resp, err := sendCommand(conn, "call_simulation 88.151.132.30,5060,9123887982,5482#+14049179360\r\n", reader)
	if err != nil {
		fmt.Println("Error during call_simulation:", err)
		os.Exit(1)
	}

	fmt.Printf("call recv: %s", call_simulation_resp)
	writeToFile("call", call_simulation_resp)

}
func sendCommand(conn net.Conn, command string, reader *bufio.Reader) (string, error) {
	// 发送命令
	fmt.Printf("Sending command: %s", command)
	_, err := conn.Write([]byte(command))
	if err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	// 设置超时时间（5秒）
	timeoutDuration := 15 * time.Second
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	// 读取多行服务器响应，拼接成一个字符串
	var responseBuilder strings.Builder
	for {
		select {
		case <-timer.C: // 超时
			log.Println("i/o timeout")
			return "", fmt.Errorf("read timeout after %v", timeoutDuration)
		default:
			// 从连接读取数据
			line, err := reader.ReadString('\n')
			if err != nil {
				// 处理读取错误
				if strings.Contains(err.Error(), "i/o timeout") {
					log.Println("i/o timeout")
					return "", fmt.Errorf("failed to read response: %w", err)
				}
				return "", fmt.Errorf("failed to read response: %w", err)
			}

			// 将读取的行追加到响应结果
			responseBuilder.WriteString(line)

			// 如果返回某些预定义的结束标记，可以在此处判断并终止读取
			if strings.HasSuffix(line, "<Call Simulation Test progress>Done</Call Simulation Test progress>") {
				break
			}

			// 重置超时定时器
			timer.Reset(timeoutDuration)
		}
	}

	return responseBuilder.String(), nil
}
