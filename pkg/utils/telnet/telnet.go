package telnet

import (
	"fmt"
	"net"
	"time"
)

// TelnetClient 定义telnet客户端结构体
type TelnetClient struct {
	IP               string
	Port             string
	IsAuthentication bool
	conn             net.Conn
}

// NewTelnetClient 创建新的telnet客户端
func NewTelnetClient(ip string, port string) *TelnetClient {
	return &TelnetClient{
		IP:               ip,
		Port:             port,
		IsAuthentication: false,
	}
}

// Connect 建立连接
func (t *TelnetClient) Connect() error {
	address := fmt.Sprintf("%s:%s", t.IP, t.Port)
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	t.conn = conn
	return nil
}

// Login 发送登录命令
func (t *TelnetClient) Login() error {
	if t.conn == nil {
		return fmt.Errorf("未建立连接")
	}

	// 设置读写超时
	t.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	t.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 发送登录命令
	_, err := t.conn.Write([]byte("login\r\n"))
	if err != nil {
		return fmt.Errorf("发送登录命令失败: %v", err)
	}

	// 读取响应
	buffer := make([]byte, 1024)
	n, err := t.conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to load: %v", err)
	}

	// 收到响应后设置认证状态为true
	//t.IsAuthentication = true
	fmt.Printf("Login Resp: %s\n", string(buffer[:n]))
	return nil
}

// Close 关闭连接
func (t *TelnetClient) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
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

	fmt.Println("登录成功!")
}
