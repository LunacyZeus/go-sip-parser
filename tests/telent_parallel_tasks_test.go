package tests

import (
	"sip-parser/pkg/utils/telnet"
	"sync"
	"testing"
)

func TestTelentParallelTasks(t *testing.T) {
	commands := []string{
		"call_simulation 88.151.128.89,5060,12156094684,7462#12156924598",
		"call_simulation 88.151.132.30,5060,12196002708,5482#+12196882815",
		"call_simulation 87.237.87.28,5060,+16026988601,7193#16023154842",
		"call_simulation 162.212.244.106,5060,+19703418675,7356#+19706899942",
		"call_simulation 208.93.43.242,5060,+16232453808,17183389338",
		"call_simulation 64.125.111.164,5060,+12154008100,1111#12152520985",
	}
	// 创建一个 WaitGroup 来等待所有的 Goroutine 完成
	var wg sync.WaitGroup

	// 启动多个 Goroutine 进行并行任务
	for i := 0; i < len(commands); i++ {
		command := commands[i]
		wg.Add(1)

		go func(taskID int) {
			defer wg.Done()

			t.Logf("Task %s started", command)
			// 创建客户端实例
			client := telnet.NewTelnetClient("192.168.1.1", "23")

			// 建立连接
			err := client.Connect()
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			defer client.Close()

			content, err := client.CallSimulation(command)
			if err != nil {
				t.Errorf("CallSimulation->%v", err)
				return
			}
			t.Log(command, content)

			_ = client.LoginOut()

		}(i)
	}

	// 等待所有的 goroutine 完成
	wg.Wait()
}
