package aria2

import (
	"context"
	"github.com/jlb0906/micro-movie/basic/common"
	proto "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	"github.com/jlb0906/micro-movie/plugins/aria2cfg"
	"github.com/jlb0906/micro-movie/plugins/miniocfg"
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
	movieCli proto.MovieSrvService
	aria2Cli rpc.Client
)

func GetAria2(ctx context.Context) rpc.Client {
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
		movieCli.AddMovie(context.TODO(), &proto.AddReq{
			Movie: &proto.Movie{
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

	cli, conf := miniocfg.GetMinio()

	for _, e := range events {
		stat, err := aria2Cli.TellStatus(e.Gid)
		if err != nil {
			logger.Error(err)
			continue
		}

		// Upload the file
		arr := strings.Split(stat.Files[0].URIs[0].URI, "/")
		objectName := arr[len(arr)-1]
		filePath := stat.Files[0].Path
		filePath = filepath.Join(aria2cfg.GetAria2().Prefix, filePath)

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

		movieCli.UpdateMovie(context.TODO(), &proto.UpdateReq{
			Movie: &proto.Movie{
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

	if aria2Cli != nil {
		return
	}

	cfg := aria2cfg.GetAria2()
	aria2Cli, _ = rpc.New(context.TODO(), cfg.Uri, cfg.Token, time.Duration(cfg.Timeout)*time.Second, new(MovieNotifier))
	movieCli = proto.NewMovieSrvService(common.MovieSrv, client.DefaultClient)
}
