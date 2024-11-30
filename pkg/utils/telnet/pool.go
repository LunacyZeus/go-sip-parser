package telnet

import (
	"fmt"
	"sync"
)

// TelnetClientPool 定义连接池结构体
type TelnetClientPool struct {
	mu         sync.Mutex
	clients    []*TelnetClient // 用于存放 TelnetClient 实例
	maxClients int             // 最大连接数
}

// NewTelnetClientPool 创建连接池
func NewTelnetClientPool(maxClients int) *TelnetClientPool {
	return &TelnetClientPool{
		clients:    make([]*TelnetClient, 0, maxClients),
		maxClients: maxClients,
	}
}

// Get 从连接池中获取一个连接
func (p *TelnetClientPool) Get(ip, port string) (*TelnetClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.clients) > 0 {
		// 如果连接池中有空闲连接，直接返回
		client := p.clients[len(p.clients)-1]
		p.clients = p.clients[:len(p.clients)-1]
		return client, nil
	}

	// 如果连接池没有空闲连接且还没达到最大连接数，创建新连接
	if len(p.clients) < p.maxClients {
		client := NewTelnetClient(ip, port)
		err := client.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to create new telnet client: %v", err)
		}
		return client, nil
	}

	return nil, fmt.Errorf("connection pool is full")
}

// Put 将连接放回连接池
func (p *TelnetClientPool) Put(client *TelnetClient) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !client.IsAvailable {
		return
	}

	if len(p.clients) < p.maxClients {
		p.clients = append(p.clients, client)
	}
}

// Close 关闭连接池中的所有连接
func (p *TelnetClientPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, client := range p.clients {
		client.Close()
	}
	p.clients = nil
}
