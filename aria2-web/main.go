package main

import (
	"fmt"
	"github.com/jlb0906/micro-movie/aria2-web/handler"
	"github.com/jlb0906/micro-movie/basic"
	"github.com/jlb0906/micro-movie/basic/common"
	"github.com/jlb0906/micro-movie/basic/config"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/config/source/grpc/v2"
	microzap "github.com/micro/go-plugins/logger/zap/v2"
	"net/http"
)

var (
	appName = "aria2_web"
	cfg     = &ariaCfg{}
)

type ariaCfg struct {
	common.AppCfg
}

func main() {
	logger.DefaultLogger, _ = microzap.NewLogger()
	initCfg()

	micReg := etcd.NewRegistry(registryOptions)

	// create new web service
	service := web.NewService(
		web.Name(cfg.Name),
		web.Version(cfg.Version),
		web.Registry(micReg),
	)

	// initialise service
	if err := service.Init(); err != nil {
		logger.Fatal(err)
	}

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	service.HandleFunc("/aria2/call", handler.Aria2Call)

	// run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	etcdCfg := &common.Etcd{}
	err := config.C().App("etcd", etcdCfg)
	if err != nil {
		panic(err)
	}

	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.Host, etcdCfg.Port)}
}

func initCfg() {
	source := grpc.NewSource(
		grpc.WithAddress("127.0.0.1:8600"),
		grpc.WithPath("/micro"),
	)

	basic.Init(
		config.WithSource(source),
		config.WithApp(appName))

	err := config.C().App(appName, cfg)
	if err != nil {
		panic(err)
	}

	logger.Infof("[initCfg] 配置 %+v", cfg)
}
