package Check

import (
	"go-twitter/Cores/mysql"
	"go-twitter/Model"
)

// unique_id数据查重
func CheckRepeatUniqueId(DataMap map[string]interface{}) bool {
	var CrawlContents []Model.CrawlContent

	maps := make(map[string]interface{})
	maps["unique_id"] = DataMap["id"]

	// 获取内容
	mysql.Db.Model(&Model.CrawlContent{}).Where(maps).Find(&CrawlContents)
	if len(CrawlContents) <= 0 {
		// 无重复
		return true
	}

	return false
}

// title数据查重
func CheckRepeatTitle(txtTitle string) bool {
	if txtTitle == "" {
		// 标题为空
		return false
	}

	var CrawlContents []Model.CrawlContent

	// 获取内容
	maps := make(map[string]interface{})
	maps["content_title"] = txtTitle
	mysql.Db.Model(&Model.CrawlContent{}).Where(maps).Find(&CrawlContents)
	if len(CrawlContents) <= 0 {
		// 无重复
		return true
	}

	return false
}
