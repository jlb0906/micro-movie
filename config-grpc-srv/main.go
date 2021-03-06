package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	"net"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
	pb "github.com/micro/go-plugins/config/source/grpc/proto"
	"google.golang.org/grpc"
)

var (
	apps = []string{"micro"}
)

type Service struct{}

func main() {
	// load config files
	err := loadConfigFile()
	if err != nil {
		logger.Fatal(err)
	}

	// new service
	service := grpc.NewServer()
	pb.RegisterSourceServer(service, new(Service))
	ts, err := net.Listen("tcp", ":8600")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("configServer started")
	err = service.Serve(ts)
	if err != nil {
		logger.Fatal(err)
	}
}

func (s Service) Read(ctx context.Context, req *pb.ReadRequest) (rsp *pb.ReadResponse, err error) {
	appName := parsePath(req.Path)
	switch appName {
	case "micro":
		rsp = &pb.ReadResponse{
			ChangeSet: getConfig(appName),
		}
		return
	default:
		err = fmt.Errorf("[Read] the first path is invalid")
		return
	}
}

func (s Service) Watch(req *pb.WatchRequest, server pb.Source_WatchServer) (err error) {
	appName := parsePath(req.Path)
	rsp := &pb.WatchResponse{
		ChangeSet: getConfig(appName),
	}
	if err = server.Send(rsp); err != nil {
		logger.Errorf("[Watch] watch files error，%s", err)
		return err
	}

	return
}

func loadConfigFile() (err error) {
	for _, app := range apps {
		if err := config.Load(file.NewSource(
			file.WithPath("./conf/" + app + ".yml"),
		)); err != nil {
			logger.Fatalf("[loadConfigFile] load files error，%s", err)
			return err
		}
	}

	// watch changes
	watcher, err := config.Watch()
	if err != nil {
		logger.Fatalf("[loadConfigFile] start watching files error，%s", err)
		return err
	}

	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				logger.Fatalf("[loadConfigFile] watch files error，%s", err)
				return
			}

			logger.Infof("[loadConfigFile] file change， %s", string(v.Bytes()))
		}
	}()

	return
}

func getConfig(appName string) *pb.ChangeSet {
	bytes := config.Get(appName).Bytes()

	logger.Infof("[getConfig] appName，%s", appName)
	return &pb.ChangeSet{
		Data:      bytes,
		Checksum:  fmt.Sprintf("%x", md5.Sum(bytes)),
		Format:    "yml",
		Source:    "file",
		Timestamp: time.Now().Unix()}
}

func parsePath(path string) (appName string) {
	paths := strings.Split(path, "/")

	if paths[0] == "" && len(paths) > 1 {
		return paths[1]
	}

	return paths[0]
}
