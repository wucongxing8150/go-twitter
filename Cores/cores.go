package Cores

import (
	"go-twitter/Cores/mysql"
	"go-twitter/Cores/redis"
)

func Init() {

	// 初始化数据库连接
	mysql.Init()

	// 初始化redis连接
	redis.Init()

}
