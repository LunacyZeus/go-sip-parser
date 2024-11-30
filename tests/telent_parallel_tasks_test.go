package tests

import (
	"fmt"
	"sip-parser/pkg/utils/telnet"
	"testing"
)

func HandleClient(tag string, pool *telnet.TelnetClientPool) {
	// 获取一个客户端实例
	client1, err := pool.Get("127.0.0.1", "4320")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Put(client1)

	if !client1.IsAuthentication {
		// 发送登录命令
		err = client1.Login()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(tag, "Successfully logged in!")
	} else {
		fmt.Println(tag, "no need login")
	}

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
	HandleClient("1", pool)
	HandleClient("2", pool)
	HandleClient("3", pool)
}
