package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/logger"
	"sync"

	"github.com/jlb0906/micro-movie/basic"
)

var (
	inited  bool
	mysqlDB *gorm.DB
	m       sync.RWMutex
)

func init() {
	basic.Register(initDB)
}

// initDB 初始化数据库
func initDB() {
	m.Lock()
	defer m.Unlock()

	if inited {
		logger.Warn(fmt.Sprint("[initDB] db 已经初始化过"))
		return
	}

	initMysql()

	inited = true
}

// GetDB 获取db
func Get() *gorm.DB {
	return mysqlDB
}
