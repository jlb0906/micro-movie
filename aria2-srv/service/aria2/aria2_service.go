package aria2

import (
	"context"
	"github.com/jlb0906/micro-movie/basic/common"
	proto "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	"github.com/jlb0906/micro-movie/plugins/aria2cfg"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/logger"
	"github.com/zyxar/argo/rpc"
	"strings"
	"sync"
	"time"
)

var (
	m        sync.RWMutex
	movieCli proto.MovieService
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
		movieCli.AddMovie(context.TODO(), &proto.AddRequest{
			Title:  arr[len(arr)-1],
			Url:    stat.Files[0].URIs[0].URI,
			Status: stat.Status,
		})
	}
}

func (MovieNotifier) OnDownloadPause(events []rpc.Event)    { logger.Infof("%s paused.", events) }
func (MovieNotifier) OnDownloadStop(events []rpc.Event)     { logger.Infof("%s stopped.", events) }
func (MovieNotifier) OnDownloadComplete(events []rpc.Event) { logger.Infof("%s completed.", events) }
func (MovieNotifier) OnDownloadError(events []rpc.Event)    { logger.Infof("%s error.", events) }

func (MovieNotifier) OnBtDownloadComplete(events []rpc.Event) {
	logger.Infof("bt %s completed.", events)
}

// Init 初始化库存服务层
func Init() {
	m.Lock()
	defer m.Unlock()

	if aria2Cli != nil {
		return
	}

	cfg := aria2cfg.GetAria2()
	aria2Cli, _ = rpc.New(context.TODO(), cfg.Uri, cfg.Token, time.Duration(cfg.Timeout)*time.Second, new(MovieNotifier))
	movieCli = proto.NewMovieService(common.MovieSrv, client.DefaultClient)
}
