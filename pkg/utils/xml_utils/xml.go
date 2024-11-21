package xml_utils

import (
	"fmt"
	"regexp"
	"strings"
)

// 获取标签名称
func getTagName(text string) string {
	re := regexp.MustCompile(`<(.*?)>`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func ParseXMLToNodeList(content string) ([]map[string]string, error) {
	// 按行分割
	lines := strings.Split(string(content), "\n")
	var nodes []map[string]string

	n := 0
	for n < len(lines) {
		line := lines[n]
		if strings.HasPrefix(line, "<") { // 根标签
			nodeName := getTagName(line)
			if nodeName == "" {
				n++
				continue
			}

			//fmt.Println(nodeName)
			var nodeContents []string

			// 遇到根标签开始获取全部内容
			for {
				line = lines[n]
				closeTag := fmt.Sprintf("</%s>", nodeName)

				if strings.Contains(line, closeTag) {
					startTag := fmt.Sprintf("<%s>", nodeName)
					// 检查是否为无子标签的情况
					if strings.HasPrefix(line, startTag) && strings.HasSuffix(line, closeTag) {
						// 无子标签
						parts := strings.Split(line, startTag)
						if len(parts) > 1 {
							data := strings.Split(parts[1], closeTag)[0]
							nodes = append(nodes, map[string]string{
								nodeName: fmt.Sprintf("%s%s%s", startTag, data, closeTag),
							})
						}
						break
					}

					// 有子标签
					nodes = append(nodes, map[string]string{
						nodeName: strings.Join(nodeContents, "\n"),
					})
					break
				}

				startTag := fmt.Sprintf("<%s>", nodeName)
				if !strings.HasPrefix(line, startTag) { // 收集子标签内容
					nodeContents = append(nodeContents, line)
				}
				n++

				if n >= len(lines) {
					break
				}
			}
		} else {
			//fmt.Printf("%v\n", []string{line})
		}
		n++
	}

	return nodes, nil
}
