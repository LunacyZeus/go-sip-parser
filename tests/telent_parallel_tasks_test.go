package tests

import (
	"fmt"
	"sip-parser/pkg/utils/telnet"
	"testing"
)

func HandleClient(tag string, pool *telnet.TelnetClientPool, command string) {
	// 获取一个客户端实例
	client, err := pool.Get("127.0.0.1", "4320")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Put(client)

	if !client.IsAuthentication {
		// 发送登录命令
		err = client.Login()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(tag, "Successfully logged in!")
	} else {
		fmt.Println(tag, "no need login")
	}

	content, err := client.CallSimulation(command)
	if err != nil {
		err = fmt.Errorf("CallSimulation->%v", err)
		return
	}
	fmt.Println(tag, command, len(content))

}

func TestTelentParallelTasks(t *testing.T) {
	_ = []string{
		"call_simulation 88.151.128.89,5060,12156094684,7462#12156924598\r\n",
		"call_simulation 88.151.132.30,5060,12196002708,5482#+12196882815\r\n",
		"call_simulation 87.237.87.28,5060,+16026988601,7193#16023154842\r\n",
		//"call_simulation 162.212.244.106,5060,+19703418675,7356#+19706899942",
		//"call_simulation 208.93.43.242,5060,+16232453808,17183389338",
		//"call_simulation 64.125.111.164,5060,+12154008100,1111#12152520985",
	}

	// 创建连接池实例
	pool := telnet.NewTelnetClientPool(10)
	go HandleClient("1", pool, "call_simulation 88.151.128.89,5060,12156094684,7462#12156924598\r\n")
	go HandleClient("2", pool, "call_simulation 88.151.132.30,5060,12196002708,5482#+12196882815\r\n")
	HandleClient("3", pool, "call_simulation 87.237.87.28,5060,+16026988601,7193#16023154842\r\n")
	select {}
}
