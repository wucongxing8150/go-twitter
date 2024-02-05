package goTwitter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go-twitter/Cores/mysql"
	"go-twitter/Cores/redis"
	"go-twitter/Model"
	"go-twitter/Services/Check"
	"go-twitter/Services/Oss"
	"go-twitter/Services/Response"
	"go-twitter/Services/Server"
	"go-twitter/Utils/Proxy"
	"go-twitter/Utils/Snowflake"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var CrawlKeyWordNum int

// GetTwitter 推特内容获取
func GetTwitter(keyWordId int64, nextToken string) bool {
	// 获取关键字详情
	CrawlKeyWord, req := getCrawlKeyWord(keyWordId)
	if req == false {
		// 关键字详情获取失败
		return false
	}

	var startTime string
	if CrawlKeyWord.CreatedAt == "" {
		// 获取7天前时间
		startTime = time.Unix(time.Now().Unix()-60*60*24*6, 0).Format(time.RFC3339)
	} else {
		startTime = CrawlKeyWord.CreatedAt
	}

	// 检测请求速率
	res := RateLimit()
	if res == false {
		// 无请求次数
		return false
	}

	// 构建参数
	result, req := getRequest(startTime, CrawlKeyWord, nextToken)
	if req == false {
		// 请求失败
		return false
	}

	if result == "" {
		fmt.Println("返回结果为空")
	}
	fmt.Println("返回结果")
	fmt.Println(result)

	// 处理返回数据为map结构
	ResponseMap, req := Response.HandleResponse(result)
	if req == false || ResponseMap == nil {
		fmt.Println("返回数据处理失败")
		return false
	}

	// 处理返回数据错误问题
	req = Response.HandleResponseError(ResponseMap)
	if req == false {
		fmt.Println("存在error字段")
		return false
	}

	// 处理返回结果中的data字段
	ResponseData, req := Response.HandleResponseData(ResponseMap)
	if req == false || ResponseData == nil {
		fmt.Println("data字段处理失败")
		return false
	}

	fmt.Println("返回结果中的data字段")
	str1, err1 := json.Marshal(ResponseData)
	if err1 != nil {
		// json解析失败
		fmt.Println(err1)
		return false
	}
	fmt.Println(string(str1))

	// 处理返回结果中的meta字段
	ResponseMeta, req := Response.HandleResponseMeta(ResponseMap)
	if req == false || ResponseMeta == nil {
		fmt.Println("meta字段处理失败")
		return false
	}

	// 处理返回数据中next_token问题
	nextToken = Response.HandleResponseMetaNextToken(ResponseMeta)

	// 处理返回结果中的includes-media字段
	ResponseIncludesMedia, req := Response.HandleResponseIncludesMedia(ResponseMap)
	if req == false || ResponseIncludesMedia == nil {
		// meta字段处理失败或不存在。此时应该判断next_token字段进行下一次请求
		return false
	}

	fmt.Println("返回结果中的includes-media字段")
	str2, err2 := json.Marshal(ResponseIncludesMedia)
	if err2 != nil {
		// json解析失败
		fmt.Println(err2)
		return false
	}
	fmt.Println(string(str2))

	// 处理data-media对应问题
	/*
		推特数据返回结构较特殊以下为例子
		data:[ 此处只有id和text文本
			{
				media_key:"1111"/["111","222"] 唯一值 用来对应数据
			},
			{
				media_key:"333"/["333","444"]
			}
		],
		media:[  此处为主体数据
			{
				media_key:"1111"
			},
			{
				media_key:"2222"
			}
			{
				media_key:"3333"
			},
			{
				media_key:"4444"
			}
		]
	*/
	for _, v := range ResponseData {
		// 检测unique_id数据重复性
		DataMap := v.(map[string]interface{})
		req := Check.CheckRepeatUniqueId(DataMap)
		if req == false {
			// 数据重复
			fmt.Println("数据重复")
			continue
		}

		mediaKeys, req := Response.HandleAttachmentsMediaKeys(DataMap)
		if req == false || mediaKeys == nil {
			// mediaKeys字段处理失败
			return false
		}

		var TwitterImgPaths []string // twitter图片地址
		var TwitterLang int          // twitter语言
		var TwitterVideoPath string  // twitter视频地址

		var AliyunImgPaths []string // 图片地址
		var AliyunVideoPath string  // 视频地址

		//  处理语言字段
		TwitterLang = Response.HandleLang(DataMap)

		// 处理media_keys字段
		for _, IncludesMedia := range ResponseIncludesMedia {
			IncludesMediaMap := IncludesMedia.(map[string]interface{})
			for _, mediaKey := range mediaKeys {
				if mediaKey == IncludesMediaMap["media_key"] {
					// 属于该数据
					if IncludesMediaMap["type"] == "video" {
						// 获取清晰度最高的视频地址
						TwitterVideoPath = Server.GetTwitterVideoPath(IncludesMediaMap)

						// 视频截图
						TwitterImgPaths = append(TwitterImgPaths, IncludesMediaMap["preview_image_url"].(string))

					} else if IncludesMediaMap["type"] == "photo" {
						// 图片
						TwitterImgPaths = append(TwitterImgPaths, IncludesMediaMap["url"].(string))
					} else {
						// 无视频/图片
						fmt.Println("无视频/图片")
						continue
					}
				}
			}
		}

		// 标题、内容处理
		txtTitle, txtContent := Server.HandleTitleContent(DataMap["text"].(string))
		if txtTitle == "" && txtContent == "" {
			// 标题/内容为空
			continue
		}

		// 检测标题数据重复性
		req = Check.CheckRepeatTitle(txtTitle)
		if req == false {
			// 数据重复
			fmt.Println("标题数据重复")
			continue
		}

		// oss视频上传
		AliyunVideoPath, req = Oss.VideoProcessing(TwitterVideoPath)
		if req == false {
			// 视频上传失败
			continue
		}

		// oss图片上传
		AliyunImgPaths = Oss.ImgProcessing(TwitterImgPaths)
		if AliyunImgPaths == nil {
			// 图片上传失败或为空
			continue
		}

		// twitter图片地址切片转为可入库数据。逗号分隔字符串
		contentImgs := Server.HandleAliyunImgPaths(AliyunImgPaths)

		// 入库
		snow, err := Snowflake.NewWorker(1)
		if err != nil {
			// 实例化失败
			continue
		}

		// 开启事务
		tx := mysql.Db.Begin()
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				log.Fatal(err)
			}
		}()
		if err := tx.Error; err != nil {
			tx.Rollback()
			continue
		}

		crawlContent := &Model.CrawlContent{
			ContentId:    snow.GetId(),
			HId:          CrawlKeyWord.HId,
			KeyWordId:    CrawlKeyWord.KeyWordId,
			UniqueId:     DataMap["id"].(string),
			Language:     TwitterLang,
			ContentTitle: txtTitle,
			Content:      txtContent,
			ContentImgs:  contentImgs,
			ContentVideo: AliyunVideoPath,
			Created:      int(time.Now().Unix()),
			Updated:      int(time.Now().Unix()),
		}

		if err := tx.Create(&crawlContent).Error; err != nil {
			tx.Rollback()
			continue
		}

		// 次数+1
		CrawlKeyWordNum++

		CrawlKeyWord.CreatedAt = DataMap["created_at"].(string)
		CrawlKeyWord.Num = CrawlKeyWordNum

		if err := tx.Model(&CrawlKeyWord).Updates(&CrawlKeyWord).Error; err != nil {
			tx.Rollback()
			continue
		}

		fmt.Println("执行结束")
		// 数量增加
		tx.Commit()
	}

	// 存在下一页数据
	if nextToken != "" {
		req := GetTwitter(CrawlKeyWord.KeyWordId, nextToken)
		if req == false {
			// 重复
			return false
		}
	}
	return true
}

// RateLimit 请求速率限制
func RateLimit() bool {
	result, err := redis.Rdb.Get(context.Background(), "goTwitterLimit").Int()
	if err != nil {
		// 创建10分钟生命周期
		_, err := redis.Rdb.Set(context.Background(), "goTwitterLimit", 1, time.Second*60*10).Result()
		if err != nil {
			// 创建失败
			return false
		}
	}

	if result >= 150 {
		// 超过速率限制
		return false
	}
	return true
}

// 请求
func getRequest(startTime string, CrawlKeyWord Model.CrawlKeyWord, nextToken string) (string, bool) {
	// 请求 from:majiji70631125 CrawlKeyWord.KeyWordName+
	a := CrawlKeyWord.KeyWordName
	b := url.QueryEscape(a)

	query := b + "+has:media_link+-is:retweet"
	// query := "from:majiji70631125+has:media_link+-is:retweet"

	client := Proxy.Proxy()

	var path string
	host := "https://api.twitter.com/2/tweets/search/recent"

	// 2022-08-23T08:18:46+08:00

	if nextToken != "" {
		path = "?query=" + query +
			"&start_time=" + url.QueryEscape(startTime) +
			"&expansions=attachments.media_keys" +
			"&tweet.fields=created_at" +
			"&media.fields=preview_image_url,url,duration_ms,alt_text,variants" +
			"&max_results=10" +
			// "&since_id=" + CrawlContent.UniqueId +
			"&next_token=" + nextToken
	} else {
		path = "?query=" + query +
			"&start_time=" + url.QueryEscape(startTime) +
			"&expansions=attachments.media_keys" +
			"&tweet.fields=created_at" +
			"&media.fields=preview_image_url,url,duration_ms,alt_text,variants" +
			// "&since_id=" + CrawlContent.UniqueId +
			"&max_results=10"
	}

	path = host + path
	fmt.Println(path)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Println("请求失败")
		return err.Error(), false
	}
	req.Header.Set("Authorization", "Bearer "+viper.GetString("twitter.bearerToken"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败")
		return err.Error(), false
	}

	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("返回结果解析失败")
		return err.Error(), false
	}
	return string(body), true
}

// 获取当前搜索关键词
func getCrawlKeyWord(keyWordId int64) (Model.CrawlKeyWord, bool) {
	var CrawlKeyWord Model.CrawlKeyWord

	maps := make(map[string]interface{})
	maps["key_word_id"] = keyWordId
	// 获取爬虫配置
	result := mysql.Db.Model(&Model.CrawlKeyWord{}).Where(maps).First(&CrawlKeyWord)
	if result.Error != nil {
		// 查询失败
		return CrawlKeyWord, false
	}

	return CrawlKeyWord, true
}

// 获取最后一条搜索内容时间
func getLastTimeCrawlContent(keyWordId int64) (Model.CrawlContent, bool) {
	var CrawlContent Model.CrawlContent

	maps := make(map[string]interface{})
	maps["key_word_id"] = keyWordId
	// 获取爬虫配置
	result := mysql.Db.Model(&Model.CrawlContent{}).Where(maps).Order("unique_id desc").First(&CrawlContent)
	if result.Error != nil {
		// 查询失败
		return CrawlContent, false
	}

	return CrawlContent, true
}
