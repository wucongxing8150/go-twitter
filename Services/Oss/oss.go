package Oss

import (
	"fmt"
	"go-twitter/Utils/aliyun"
	"strconv"
	"time"
)

// oss视频上传
func VideoProcessing(TwitterVideoPath string) (string, bool) {
	var AliyunVideoPath string
	fmt.Println(TwitterVideoPath)

	if TwitterVideoPath == "" {
		// 可能存在无视频
		return "", true
	}

	// 上传
	myVideoObject := fmt.Sprintf("%s%s%s%s%s%s%s%s", strconv.Itoa(time.Now().Year()), strconv.Itoa(int(time.Now().Month())), strconv.Itoa(time.Now().Day()), strconv.Itoa(time.Now().Hour()), strconv.Itoa(time.Now().Minute()), strconv.Itoa(time.Now().Second()), strconv.Itoa(int(time.Now().Unix())), ".mp4")

	AliyunVideoPath, err := aliyun.UploadHttpFile(myVideoObject, TwitterVideoPath)
	if err != nil {
		// 上传失败
		return "", false
	}

	return AliyunVideoPath, true
}

// ImgProcessing 图片上传
func ImgProcessing(TwitterImgPaths []string) []string {
	fmt.Println(TwitterImgPaths)
	var AliyunImgPaths []string // 图片地址
	if TwitterImgPaths == nil {
		return nil
	}

	for _, path := range TwitterImgPaths {
		// 上传
		myImgObject := fmt.Sprintf("%s%s%s%s%s%s%s%s", strconv.Itoa(time.Now().Year()), strconv.Itoa(int(time.Now().Month())), strconv.Itoa(time.Now().Day()), strconv.Itoa(time.Now().Hour()), strconv.Itoa(time.Now().Minute()), strconv.Itoa(time.Now().Second()), strconv.Itoa(int(time.Now().Unix())), ".jpg")

		AliyunImgPath, err := aliyun.UploadHttpFile(myImgObject, path)
		if err != nil {
			// 上传失败
			return nil
		}

		AliyunImgPaths = append(AliyunImgPaths, AliyunImgPath)
	}

	if AliyunImgPaths == nil {
		// 全部上传失败
		return nil
	}

	return AliyunImgPaths
}
