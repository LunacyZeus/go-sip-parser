package utils

import (
	"strings"
)

func GetPhonePart(input string) string {
	// 去除 <sip: 和 @ 后面的内容
	if strings.Contains(input, "<sip:") {
		// 截取 <sip: 后的部分，去除后面的 @ 和 IP
		start := strings.Index(input, "<sip:") + len("<sip:")
		end := strings.Index(input[start:], "@")
		if end != -1 {
			phone := input[start : start+end]
			// 打印号码部分
			if strings.Contains(phone, "#") {
				tmp := strings.Split(phone, "#")
				phone = tmp[1]
			}

			if len(phone) > 0 && phone[0] == '+' {
				phone = phone[1:]
			}

			return phone
		}
	} else {
		// 如果不是 <sip: 格式，按 @ 分割并打印号码部分
		parts := strings.Split(input, "@")
		if len(parts) > 0 {
			return parts[0] // 打印号码部分
		}
	}
	return ""
}
