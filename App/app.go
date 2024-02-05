package App

import (
	"fmt"
	"go-twitter/App/goTwitter"
	"go-twitter/Cores/mysql"
	"go-twitter/Model"
	"log"
	"time"
)

func Init() {
	for {
		// 获取爬虫配置
		CrawlSettings, res := getCrawlSettings()
		if res == false || CrawlSettings == nil {
			log.Println("无爬取任务")
			time.Sleep(time.Second * 60 * 3) // 睡眠10分钟
			continue
		}

		// 获取全部搜索关键词
		CrawlKeyWords, res := getCrawlKeyWords()
		if res == false || CrawlKeyWords == nil {
			log.Println("无关键字")
			continue
		}

		// 处理业务
		for _, v := range CrawlKeyWords {
			if v.SettingPlatform == 1 {
				fmt.Println(v)
				req := goTwitter.GetTwitter(v.KeyWordId, "")
				if req == false {
					// 重复/失败
					continue
				}
			} else if v.SettingPlatform == 2 {
				log.Println("facebook")
				continue
			} else {
				log.Println("未知")
				continue
			}
		}

		// 全部搜索完毕，此处暂停10分钟，防止浪费性能
		time.Sleep(time.Second * 60 * 3)
	}

}

// 获取爬虫配置
func getCrawlSettings() ([]Model.CrawlSetting, bool) {
	var CrawlSettings []Model.CrawlSetting

	maps := make(map[string]interface{})

	// 获取爬虫配置
	mysql.Db.Model(&Model.CrawlSetting{}).Where(maps).Find(&CrawlSettings)
	if len(CrawlSettings) <= 0 {
		log.Println("无爬取任务")
		return nil, false
	}

	return CrawlSettings, true
}

// 获取全部搜索关键词
func getCrawlKeyWords() ([]Model.CrawlKeyWord, bool) {
	var CrawlKeyWords []Model.CrawlKeyWord
	maps := make(map[string]interface{})

	// 获取爬虫配置
	mysql.Db.Model(&Model.CrawlKeyWord{}).Where(maps).Find(&CrawlKeyWords)
	if len(CrawlKeyWords) <= 0 {
		log.Println("无关键字")
		return nil, false
	}

	return CrawlKeyWords, true
}
