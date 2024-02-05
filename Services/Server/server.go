package Server

import (
	"go-twitter/Utils/Pcre"
	"strings"
)

// 获取清晰度最高的视频地址
func GetTwitterVideoPath(IncludesMediaMap map[string]interface{}) string {
	var TwitterVideoPath string
	var bitRate float64 // 视频大小字段

	// 寻找格式最大的视频
	variants := IncludesMediaMap["variants"].([]interface{})
	for _, variant := range variants {
		variantMap := variant.(map[string]interface{})
		if variantMap["bit_rate"] != nil {
			if variantMap["bit_rate"].(float64) > bitRate {
				// 视频地址
				TwitterVideoPath = variantMap["url"].(string)
				bitRate = variantMap["bit_rate"].(float64)
			}
		}
	}

	return TwitterVideoPath
}

// 标题、内容处理
func HandleTitleContent(text string) (string, string) {
	var txtTitle string   // 标题
	var txtContent string // 内容

	if text == "" {
		// 标题为空
		return txtTitle, txtContent
	}
	// 正则处理所有twitter网址，赋空
	txtTitle = Pcre.PergTwitterHttp(text)

	// 匹配title 如存在.第一句赋为标题，后续设为内容
	comma := strings.Index(txtTitle, ".")
	if comma != -1 {
		txtContent = txtTitle[comma+1:]
		txtTitle = txtTitle[:comma]
	}

	// 去除换行符
	// txtTitle = strings.Replace(txtTitle, "\n", "", -1)
	// txtContent = strings.Replace(txtContent, "\n", "", -1)

	return txtTitle, txtContent
}

// 处理阿里云图片，切片转为字符串
func HandleAliyunImgPaths(AliyunImgPaths []string) string {
	var result string
	if len(AliyunImgPaths) > 0 {
		// 图片
		for i, i2 := range AliyunImgPaths {
			if i == 0 {
				result = i2
			} else {
				result = result + "," + i2
			}
		}
	}

	return result
}
