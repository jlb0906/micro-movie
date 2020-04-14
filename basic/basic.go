package basic

import (
	"github.com/jlb0906/micro-movie/basic/config"
)

var (
	pluginFuncs []func()
)

func Init(opts ...config.Option) {
	// 初始化配置
	config.Init(opts...)

	// 加载依赖配置的插件
	for _, f := range pluginFuncs {
		f()
	}
}

func Register(f func()) {
	pluginFuncs = append(pluginFuncs, f)
}
