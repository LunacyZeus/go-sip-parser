package process

import (
	"io/ioutil"
	"log"
	"sip-parser/pkg/utils/rate"
)

func TestFunc() {
	// 指定要读取的文件路径
	filePath := "call.xml"

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	// 将读取的内容转换为字符串
	content := string(data)

	rate.ParseRateFromContent("", content)

}
