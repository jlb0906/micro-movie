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

	var err error

	if inited {
		err = fmt.Errorf("[initDB] db 已经初始化过")
		logger.Error(err.Error())
		return
	}

	initMysql()

	inited = true
}

// GetDB 获取db
func GetDB() *gorm.DB {
	return mysqlDB
}
