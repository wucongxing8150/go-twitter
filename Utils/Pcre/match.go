// Package Pcre 正则匹配
package Pcre

import (
	"regexp"
	"strings"
)

const (
	pcreTwitterHttp = "https://t.co/[0-9a-zA-Z]{10}" // 例 https://t.co/4dTnI0jUSt
)

// 匹配内容中的推特网址
func PergTwitterHttp(content string) string {
	compile := regexp.MustCompile(pcreTwitterHttp)
	submatch := compile.FindAllStringSubmatch(content, -1)
	for _, text := range submatch {
		content = strings.Replace(content, text[0], "", 1) // 替换匹配到的字符串为空
	}
	return content
}
