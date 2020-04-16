package main

import (
	"fmt"
	"github.com/jlb0906/micro-movie/basic"
	"github.com/jlb0906/micro-movie/basic/common"
	"github.com/jlb0906/micro-movie/basic/config"
	"github.com/jlb0906/micro-movie/movie-srv/handler"
	pb "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	"github.com/jlb0906/micro-movie/movie-srv/service"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/config/source/grpc/v2"
	microzap "github.com/micro/go-plugins/logger/zap/v2"
)

var (
	appName = "movie_srv"
	cfg     = &movieCfg{}
)

type movieCfg struct {
	common.AppCfg
}

func main() {
	logger.DefaultLogger, _ = microzap.NewLogger()

	// 初始化配置、数据库等信息
	initCfg()

	// 使用etcd注册
	micReg := etcd.NewRegistry(registryOptions)

	// New Service
	srv := micro.NewService(
		micro.Name(cfg.Name),
		micro.Version(cfg.Version),
		micro.Registry(micReg),
	)

	// Initialise service
	srv.Init(
		micro.Action(func(c *cli.Context) error {
			// 初始化模型层
			service.Init()
			// 初始化handler
			handler.Init()
			return nil
		}),
	)

	// Register Handler
	pb.RegisterMovieHandler(srv.Server(), new(handler.Movie))

	// Run service
	if err := srv.Run(); err != nil {
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
