package process

import (
	"fmt"
	"os"
	"sip-parser/pkg/telnet"
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
	tc := telnet.TelnetClient{
		Address:  "127.0.0.1",
		Port:     "4320",
		Login:    "",
		Password: "",
	}
	err := tc.Dial()
	if err != nil {
		fmt.Printf("failed open connect with error = %v\n", err)
		return
	}
	defer tc.Close()

	stdout, err := tc.Execute("login")
	if err != nil {
		fmt.Printf("failed execute command with error = %v\n", err)
		return
	}

	fmt.Printf(string(stdout))

}
