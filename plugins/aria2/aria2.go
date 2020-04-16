package aria2

import (
	"github.com/jlb0906/micro-movie/basic"
	"github.com/jlb0906/micro-movie/basic/config"
	"github.com/micro/go-micro/v2/logger"
	"sync"
)

var (
	c      *Conf
	m      sync.RWMutex
	inited bool
)

// 配置
type Conf struct {
	Uri     string `json:"uri"`
	Token   string `json:"token"`
	Timeout int    `json:"timeout"`
	Prefix  string `json:"prefix"`
}

// 初始化
func init() {
	basic.Register(initAria2)
}

func initAria2() {
	m.Lock()
	defer m.Unlock()

	if inited {
		logger.Warn("[initAria2] 已经初始化过Aria2...")
		return
	}

	logger.Infof("[initAria2] 初始化Aria2...")

	src := config.C()
	c = new(Conf)
	err := src.Path("aria2", c)
	if err != nil {
		logger.Fatal("[initAria2] %s", err)
	}

	inited = true
	logger.Infof("[initAria2] 成功")
}

func Get() *Conf {
	return c
}
