package aria2

import (
	"context"
	"fmt"
	"github.com/jlb0906/micro-movie/basic/common"
	pb "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	aria2cfg "github.com/jlb0906/micro-movie/plugins/aria2"
	miniocfg "github.com/jlb0906/micro-movie/plugins/minio"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/logger"
	"github.com/minio/minio-go/v6"
	"github.com/zyxar/argo/rpc"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	m        sync.RWMutex
	inited   bool
	movieCli pb.MovieService
	aria2Cli rpc.Client
)

func GetAria2() rpc.Client {
	return aria2Cli
}

type MovieNotifier struct{}

func (MovieNotifier) OnDownloadStart(events []rpc.Event) {
	logger.Infof("%s started.", events)
	for _, e := range events {
		stat, err := aria2Cli.TellStatus(e.Gid)
		if err != nil {
			logger.Error(err)
			continue
		}
		arr := strings.Split(stat.Files[0].URIs[0].URI, "/")
		movieCli.AddMovie(context.TODO(), &pb.AddReq{
			Movie: &pb.MovieMsg{
				Id:     e.Gid,
				Title:  arr[len(arr)-1],
				Url:    stat.Files[0].URIs[0].URI,
				Status: stat.Status,
			},
		})
	}
}

func (MovieNotifier) OnDownloadPause(events []rpc.Event) { logger.Infof("%s paused.", events) }
func (MovieNotifier) OnDownloadStop(events []rpc.Event)  { logger.Infof("%s stopped.", events) }
func (MovieNotifier) OnDownloadComplete(events []rpc.Event) {
	logger.Infof("%s completed.", events)

	cli, conf := miniocfg.Get()

	for _, e := range events {
		stat, err := aria2Cli.TellStatus(e.Gid)
		if err != nil {
			logger.Error(err)
			continue
		}

		// Upload the file
		filePath := stat.Files[0].Path
		arr := strings.Split(filePath, "/")
		objectName := arr[len(arr)-1]
		filePath = filepath.Join(aria2cfg.Get().Prefix, filePath)

		// Upload the file with FPutObject
		n, err := cli.FPutObject(conf.BucketName, objectName, filePath, minio.PutObjectOptions{})
		if err != nil {
			logger.Error(err)
			return
		}

		// 生成7天访问地址
		u, err := cli.PresignedGetObject(conf.BucketName, objectName, 604800*time.Second, map[string][]string{})
		if err != nil {
			logger.Error(err)
			return
		}

		movieCli.UpdateMovie(context.TODO(), &pb.UpdateReq{
			Movie: &pb.MovieMsg{
				Id:     e.Gid,
				Title:  objectName,
				Url:    u.String(),
				Status: stat.Status,
			},
		})

		logger.Infof("Successfully uploaded %s of size %d", objectName, n)
	}
}
func (MovieNotifier) OnDownloadError(events []rpc.Event) { logger.Infof("%s error.", events) }

func (MovieNotifier) OnBtDownloadComplete(events []rpc.Event) {
	logger.Infof("bt %s completed.", events)
}

// Init 初始化服务层
func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		logger.Warn(fmt.Sprint("[Init] service 已经初始化过"))
		return
	}

	cfg := aria2cfg.Get()
	aria2Cli, _ = rpc.New(context.TODO(), cfg.Uri, cfg.Token, time.Duration(cfg.Timeout)*time.Second, new(MovieNotifier))
	movieCli = pb.NewMovieService(common.MovieSrv, client.DefaultClient)
	inited = true
}
