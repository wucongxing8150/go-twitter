// @Description: mysql初始化连接
// @Author: wucongxing
// @Date:2021/12/24 14:00

package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"go-twitter/Config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`                           // 服务器地址
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`                           // 端口
	DbName       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                  // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`               // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`               // 数据库密码
	MaxIdleConns int    `mapstructure:"max-idle-cons" json:"MaxIdleConns" yaml:"max-idle-cons"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-cons" json:"MaxOpenConns" yaml:"max-open-cons"` // 打开到数据库的最大连接数
	Debug        bool   `mapstructure:"debug" json:"debug" yaml:"debug"`                        // 是否开启Gorm全局日志
}

var Db *gorm.DB
var m *Config.Mysql

func Init() *gorm.DB {
	var err error

	if viper.GetString("env") == "product" {
		m = Config.GetMasterMysql()
	} else {
		m = Config.GetDevMysql()
	}

	// var m Mysql
	// m.Username = "userName"
	// m.Password = "password"
	// m.Host = "host"
	// m.Port = 3306
	// m.DbName = "DbName"

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", m.Username,
		m.Password, m.Host, m.Port, m.DbName, "10s")
	fmt.Println(dsn)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	sqlDB, _ := Db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 调试模式
	// Db.LogMode(m.Debug)    // 打印sql
	// Db.SingularTable(true) // 全局禁用表名复数
	fmt.Println("数据库连接成功......")
	return Db
}
