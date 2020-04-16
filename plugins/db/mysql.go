package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jlb0906/micro-movie/basic/config"
	"github.com/micro/go-micro/v2/logger"
)

type db struct {
	Mysql mysqlConf `json:"mysql"`
}

// Mysql 配置
type mysqlConf struct {
	URL    string `json:"url"`
	Enable bool   `json:"enabled"`
}

func initMysql() {
	logger.Infof("[initMysql] 初始化Mysql")

	c := config.C()
	cfg := &db{}

	err := c.App("db", cfg)
	if err != nil {
		logger.Fatalf("[initMysql] %s", err)
	}
	logger.Infof("mysql配置信息: %+v", cfg)

	if !cfg.Mysql.Enable {
		logger.Infof("[initMysql] 未启用Mysql")
		return
	}

	// 创建连接
	mysqlDB, err = gorm.Open("mysql", cfg.Mysql.URL)
	if err != nil {
		logger.Fatal(err)
	}

	if mysqlDB != nil {
		mysqlDB.SingularTable(true)
	}

	logger.Infof("[initMysql] Mysql 连接成功")
}
