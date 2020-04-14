package aria2cfg

import (
	"github.com/jlb0906/micro-movie/basic"
	"github.com/jlb0906/micro-movie/basic/config"
	"github.com/micro/go-micro/v2/logger"
	"sync"
)

var (
	cfg    *Aria2
	m      sync.RWMutex
	inited bool
)

// Aria2 Aria2 配置
type Aria2 struct {
	Uri     string `json:"uri"`
	Token   string `json:"token"`
	Timeout int    `json:"timeout"`
}

// init 初始化Redis
func init() {
	basic.Register(initAria2)
}

func initAria2() {
	m.Lock()
	defer m.Unlock()

	if inited {
		logger.Infof("[initAria2] 已经初始化过Aria2...")
		return
	}

	logger.Infof("[initAria2] 初始化Aria2...")

	c := config.C()
	cfg = new(Aria2)
	err := c.Path("aria2", cfg)
	if err != nil {
		logger.Error("[initAria2] %s", err)
	}

	if err != nil {
		logger.Error(err)
	}

	logger.Infof("[initAria2] Arias，成功")
}

// Aria2 获取aria2
func GetAria2() *Aria2 {
	return cfg
}
