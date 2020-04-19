package aria2

import (
	"context"
	"fmt"
	"github.com/jlb0906/micro-movie/basic/common"
	pb "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	aria2cfg "github.com/jlb0906/micro-movie/plugins/aria2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/logger"
	"github.com/zyxar/argo/rpc"
	"strings"
	"sync"
	"time"
)

var (
	m        sync.RWMutex
	inited   bool
	movieCli pb.MovieService
	aria2Cli rpc.Client
	e        Engine
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
		movieCli.UpdateMovieByGid(context.TODO(), &pb.UpdateReq{
			Movie: &pb.MovieMsg{
				Gid:    e.Gid,
				Title:  arr[len(arr)-1],
				Uri:    stat.Files[0].URIs[0].URI,
				Status: stat.Status,
			},
		})
	}
}

func (MovieNotifier) OnDownloadPause(events []rpc.Event) { logger.Infof("%s paused.", events) }
func (MovieNotifier) OnDownloadStop(events []rpc.Event)  { logger.Infof("%s stopped.", events) }
func (MovieNotifier) OnDownloadComplete(events []rpc.Event) {
	logger.Infof("%s completed.", events)
	e.Submit(events)
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

	e = NewAsyncEngine()
	e.Init(cfg.WorkerCount)

	inited = true
}
