package main

import (
	"go-twitter/App"
	"go-twitter/Cores"
	vip "go-twitter/Cores/viper"
	"time"
)

func main() {
	// 初始化Viper 加载配置
	vip.Init()

	// 加载核心方法
	Cores.Init()

	// 业务处理
	App.Init()

}

func ParseTime(layout string, timeStr string) (time.Time, error) {
	var cstZone = time.FixedZone("CST", 8*3600)
	return time.ParseInLocation(layout, timeStr, cstZone)
}
