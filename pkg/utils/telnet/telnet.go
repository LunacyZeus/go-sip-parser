package telnet

import (
	"bufio"
	"fmt"
	"github.com/google/uuid" // 引入uuid库
	"log"
	"net"
	"strings"
	"time"
)

// TelnetClient 定义telnet客户端结构体
type TelnetClient struct {
	UUID             string
	IP               string
	Port             string
	IsAuthentication bool
	IsAvailable      bool
	conn             net.Conn
}

// NewTelnetClient 创建新的telnet客户端
func NewTelnetClient(ip string, port string) *TelnetClient {
	return &TelnetClient{
		UUID:             uuid.NewString(),
		IP:               ip,
		Port:             port,
		IsAuthentication: false,
		IsAvailable:      false,
	}
}

// Connect 建立连接
func (t *TelnetClient) Connect() error {
	address := fmt.Sprintf("%s:%s", t.IP, t.Port)
	conn, err := net.DialTimeout("tcp", address, 15*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	t.conn = conn
	t.IsAvailable = true
	return nil
}

// Login 发送登录命令
func (t *TelnetClient) Login() error {
	if t.conn == nil {
		return fmt.Errorf("cannot establish connection")
	}

	// 设置读写超时
	//t.conn.SetReadDeadline(time.Now().Add(45 * time.Second))
	//t.conn.SetWriteDeadline(time.Now().Add(45 * time.Second))

	// 发送登录命令
	_, err := t.conn.Write([]byte("login\r\n"))
	if err != nil {
		return fmt.Errorf("Failed to send command: %v", err)
	}

	// 读取响应
	buffer := make([]byte, 1024)
	_, err = t.conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to load: %v", err)
	}

	// 收到响应后设置认证状态为true
	t.IsAuthentication = true
	//log.Printf("Login Resp: %s\n", string(buffer[:n]))
	return nil
}

// Ping
func (t *TelnetClient) Ping() error {
	if t.conn == nil {
		return fmt.Errorf("cannot establish connection")
	}

	// 发送获取状态命令
	_, err := t.conn.Write([]byte("get_running_status\r\n"))
	if err != nil {
		return fmt.Errorf("Failed to send command: %v", err)
	}

	// 读取响应
	buffer := make([]byte, 1024)
	n, err := t.conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to load: %v", err)
	}

	// 收到响应后设置认证状态为true
	//t.IsAuthentication = true
	log.Printf("running Resp: %s\n", string(buffer[:n]))
	return nil
}

// LoginOut 发送登录命令
func (t *TelnetClient) LoginOut() error {
	if t.conn == nil {
		return fmt.Errorf("cannot establish connection")
	}

	// 设置读写超时
	//t.conn.SetReadDeadline(time.Now().Add(45 * time.Second))
	//t.conn.SetWriteDeadline(time.Now().Add(45 * time.Second))

	// 发送注销命令
	_, err := t.conn.Write([]byte("logout\r\n"))
	if err != nil {
		return fmt.Errorf("Failed to send command: %v", err)
	}

	// 读取响应
	buffer := make([]byte, 1024)
	_, err = t.conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to load: %v", err)
	}

	// 收到响应后设置认证状态为true
	//t.IsAuthentication = true
	//fmt.Printf("LoginOut Resp: %s\n", string(buffer[:n]))
	t.IsAvailable = false
	return nil
}

// Close 关闭连接
func (t *TelnetClient) Close() (err error) {
	t.IsAvailable = false
	if t.conn != nil {
		err = t.conn.Close()
	}
	return
}

func (t *TelnetClient) CloseLogout() (err error) {
	t.IsAvailable = false
	if t.conn != nil {
		//t.LoginOut()
		err = t.conn.Close()
	}
	return
}

// CallSimulation 发送 call_simulation 命令并读取完整响应
func (t *TelnetClient) CallSimulation(command string) (string, error) {
	if t.conn == nil {
		return "", fmt.Errorf("cannot establish connection")
	}

	// 构建命令
	//command := fmt.Sprintf("call_simulation %s,%s,%s,%s\r\n", callerIp, callerPort, ani, dnis)

	// 发送命令
	_, err := t.conn.Write([]byte(command))
	if err != nil {
		return "", fmt.Errorf("Failed to send command: %v", err)
	}

	// 读取完整响应
	reader := bufio.NewReaderSize(t.conn, 8192)
	var response strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		response.WriteString(line)
		// 假设服务器返回的响应以某个特定标识符结束，比如 "END"
		if strings.Contains(line, "<Call Simulation Test progress>Done</Call Simulation Test progress>") {
			break
		}
		if strings.Contains(line, "</Origination-State>") {
			break
		}

	}

	//fmt.Printf("Call Simulation Resp: %s\n", response.String())
	return response.String(), nil
}

// 使用示例
func main() {
	// 创建客户端实例
	client := NewTelnetClient("192.168.1.1", "23")

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

	fmt.Println("successfully login!")
}
