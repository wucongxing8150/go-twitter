package aliyun

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"go-twitter/Utils/Proxy"
	"net/http"
)

// UploadLocalFile 上传文件
func UploadLocalFile(myObject string, localPath string) (url string, err error) {
	Endpoint := viper.GetString("aliyun.Endpoint")
	accessKey := viper.GetString("aliyun.accessKey")
	accessSecret := viper.GetString("aliyun.accessSecret")
	client, err := oss.New(Endpoint, accessKey, accessSecret)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(viper.GetString("aliyun.Bucket"))
	if err != nil {
		return "", err
	}

	myObject = "content/" + myObject // 文件名称-路径
	err = bucket.PutObjectFromFile(myObject, localPath)
	if err != nil {
		return "", err
	}

	// 拼接url
	path := "https://" + viper.GetString("aliyun.Bucket") + "." + Endpoint + "/" + myObject

	return path, nil
}

// UploadHttpFile 上传http文件
func UploadHttpFile(myObject string, httpPath string) (url string, err error) {
	Endpoint := viper.GetString("aliyun.Endpoint")
	accessKey := viper.GetString("aliyun.accessKey")
	accessSecret := viper.GetString("aliyun.accessSecret")
	client, err := oss.New(Endpoint, accessKey, accessSecret)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(viper.GetString("aliyun.Bucket"))
	if err != nil {
		return "", err
	}

	httpClient := Proxy.Proxy()

	req, err := http.NewRequest("GET", httpPath, nil)
	if err != nil {
		fmt.Println("上传失败")
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("上传失败")
		return "", err
	}

	defer resp.Body.Close()

	myObject = "content/" + myObject // 文件名称-路径
	err = bucket.PutObject(myObject, resp.Body)
	if err != nil {
		return "", err
	}

	// 拼接url
	path := "https://" + viper.GetString("aliyun.Bucket") + "." + Endpoint + "/" + myObject

	return path, nil
}

// VideoScreenshot 阿里云oss视频截图
func VideoScreenshot(aliyunFilePath string) string {

	// 拼接url
	path := aliyunFilePath + "?x-oss-process=video/snapshot,t_3000,f_png,w_0,h_0,m_fast,ar_auto"

	return path
}
