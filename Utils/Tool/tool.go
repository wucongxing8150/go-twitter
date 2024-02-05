package Tool

import (
	"strconv"
)

// 判断map[key]是否存在
func MapKeyExist(mapValue map[string]interface{}, key string) interface{} {
	if value, ok := mapValue[key]; ok {
		return value
	} else {
		return nil
	}
}

// 过滤 emoji 表情
func FilterEmoji(content string) string {
	runes := []rune(content)
	res := ""

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r < 128 {
			res += string(r)
		} else {
			res += "&#" + strconv.FormatInt(int64(r), 10) + ";"
		}
	}
	return res
}
