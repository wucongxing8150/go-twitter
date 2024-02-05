package Response

import (
	"encoding/json"
	"fmt"
	"go-twitter/Utils/Tool"
)

// 处理返回数据为map结构
func HandleResponse(result string) (map[string]interface{}, bool) {
	var ResponseMap map[string]interface{}

	err := json.Unmarshal([]byte(result), &ResponseMap)
	if err != nil {
		// json解析失败
		fmt.Println(err)
		return ResponseMap, false
	}

	return ResponseMap, true
}

// 处理返回数据错误问题
func HandleResponseError(ResponseMap map[string]interface{}) bool {
	// 判断error是否存在
	errorVal := Tool.MapKeyExist(ResponseMap, "error")
	if errorVal != nil {
		// 存在错误
		fmt.Println(errorVal)
		return false
	}
	return true
}

// 处理返回结果中的data字段
func HandleResponseData(ResponseMap map[string]interface{}) ([]interface{}, bool) {
	var ResponseData []interface{}

	dataVal := Tool.MapKeyExist(ResponseMap, "data")
	if dataVal == nil {
		// meta不存在
		fmt.Println(dataVal)
		return ResponseData, false
	}
	ResponseData = dataVal.([]interface{})
	return ResponseData, true
}

// 处理返回结果中的meta字段
func HandleResponseMeta(ResponseMap map[string]interface{}) (interface{}, bool) {
	var ResponseMeta interface{}

	// 判断meta是否存在 内有next_token
	ResponseMeta = Tool.MapKeyExist(ResponseMap, "meta")
	if ResponseMeta == nil {
		// meta不存在
		fmt.Println(ResponseMeta)
		return ResponseMeta, false
	}

	return ResponseMeta, true
}

// 处理返回数据中next_token问题
func HandleResponseMetaNextToken(ResponseMeta interface{}) string {
	// 判断meta是否存在 内有next_token
	meta := ResponseMeta.(map[string]interface{})
	if meta["next_token"] != nil {
		return meta["next_token"].(string)
	}

	return ""
}

// 处理返回结果中的includes.media字段
func HandleResponseIncludesMedia(ResponseMap map[string]interface{}) ([]interface{}, bool) {

	var ResponseIncludesMedia []interface{}

	// 判断includes是否存在 主题数据
	includesVal := Tool.MapKeyExist(ResponseMap, "includes")
	if includesVal == nil {
		// includes不存在
		return nil, false
	}

	// 判断media是否存在 media字段数据
	mediaVal := Tool.MapKeyExist(includesVal.(map[string]interface{}), "media")
	if mediaVal == nil {
		// media不存在
		return nil, false
	}

	ResponseIncludesMedia = mediaVal.([]interface{})

	return ResponseIncludesMedia, true
}

// 处理返回结果中的ResponseData-Attachments-MediaKey数据
func HandleAttachmentsMediaKeys(DataMap map[string]interface{}) ([]interface{}, bool) {
	if DataMap["attachments"] == nil {
		// 无视频图片，资源不符合要求
		fmt.Println("无视频图片，资源不符合要求")
		return nil, false
	}

	attachments := DataMap["attachments"].(map[string]interface{})
	mediaKeys := attachments["media_keys"].([]interface{}) // data中的media_key字段，切片形式，当一个推文包含多个图片时，此字段为切片
	return mediaKeys, true
}

// 处理返回结果中的ResponseData-lang语言字段
func HandleLang(DataMap map[string]interface{}) int {
	if DataMap["lang"] == nil {
		return 0
	}

	if DataMap["lang"] == "en" {
		return 1
	} else if DataMap["lang"] == "ar" {
		return 2
	}
	return 0
}
