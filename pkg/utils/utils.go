package utils

import (
	"regexp"
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

// 判断 SIP 地址是否为呼出或呼入
func IsOutbound(sip string) bool {
	// 使用正则表达式提取 SIP 地址中的 IP 部分
	re := regexp.MustCompile(`@([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
	matches := re.FindStringSubmatch(sip)

	//log.Println("matches:", matches)
	if len(matches) > 1 {
		// 提取到的 IP 地址
		ip := matches[1]

		if len(ip) > 0 && ip[0] == '@' {
			ip = ip[1:]
		}
		//fmt.Println(ip)

		// 判断 IP 是否以 "172." 开头，判断是否是呼出
		if strings.HasPrefix(ip, "172.") {
			return true
		}
	}
	// 如果没有匹配到或者 IP 不是以 "172." 开头，则返回 false，表示是呼入
	return false
}
