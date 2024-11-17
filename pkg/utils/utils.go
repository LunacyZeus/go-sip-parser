package utils

import (
	"strings"
	"unicode"
)

// 定义 SIP 请求方法列表
var sipRequestMethods = []string{
	"INVITE", "ACK", "CANCEL", "BYE", "INFO", "PRACK", "UPDATE", "OPTIONS", "MESSAGE",
}

// cleanLine 清理行中的非打印字符
func cleanLine(line string) string {
	var result []rune
	for _, r := range line {
		if unicode.IsPrint(r) {
			result = append(result, r)
		}
	}
	return string(result)
}

// 获取请求行
func GetRequestLine(line string) (string, string) {
	// 清理行中的非打印字符
	cleanedLine := cleanLine(line)

	// 遍历请求方法列表，检查是否包含有效的请求方法
	for _, method := range sipRequestMethods {
		if strings.Contains(cleanedLine, method) {
			tmp := strings.SplitN(cleanedLine, method+" ", 2)
			if len(tmp) > 1 {
				return method, method + " " + tmp[1]
			}
		}
	}

	// 如果包含 SIP/2.0 响应行，提取响应部分
	if strings.Contains(cleanedLine, "SIP/2.0 ") {
		tmp := strings.SplitN(cleanedLine, "SIP/2.0 ", 2)
		if len(tmp) > 1 {
			return "", "SIP/2.0 " + tmp[1]
		}
	}

	// 如果没有匹配到任何请求行或响应行，则返回空字符串
	return "", ""
}
