package handler

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/jlb0906/micro-movie/aria2-srv/proto/aria2"
	aria2srv "github.com/jlb0906/micro-movie/aria2-srv/service/aria2"
	"github.com/jlb0906/micro-movie/basic/common"
	"github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/zyxar/argo/rpc"
	"sync"
)

var (
	m        sync.RWMutex
	inited   bool
	aria2Cli rpc.Client
	movieSrv movie.MovieService
)

type Aria2 struct{}

func (e *Aria2) AddURI(ctx context.Context, req *pb.AddURIReq, rsp *pb.AddURIRsp) error {
	logger.Info("Received Aria2.AddURI request")

	gid, err := aria2Cli.AddURI(req.Uri)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &pb.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("添加了下载任务：%v", gid)

	rsp.Gid = gid
	return nil
}

func (e *Aria2) Remove(ctx context.Context, req *pb.RemoveReq, rsp *pb.RemoveRsp) error {
	logger.Info("Received Aria2.Remove request")

	gid, err := aria2Cli.Remove(req.Gid)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &pb.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("删除了下载任务：%v", gid)

	rsp.Gid = "gid"
	return nil
}

func (e *Aria2) Pause(ctx context.Context, req *pb.PauseReq, rsp *pb.PauseRsp) error {
	logger.Info("Received Aria2.Pause request")

	gid, err := aria2Cli.Pause(req.Gid)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &pb.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("暂停了下载任务：%v", gid)

	rsp.Gid = "gid"
	return nil
}

func (e *Aria2) TellStatus(ctx context.Context, req *pb.TellStatusReq, rsp *pb.TellStatusRsp) error {
	logger.Info("Received Aria2.TellStatus request")

	info, err := aria2Cli.TellStatus(req.Gid, req.Keys...)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &pb.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("下载任务的状态：%v", info)

	data, _ := json.Marshal(info)
	rsp.Info = &pb.StatusInfo{}
	err = json.Unmarshal(data, rsp.Info)
	if err != nil {
		logger.Error(err)
	}
	return nil
}

func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		logger.Warn(fmt.Sprint("[Init] handler 已经初始化过"))
		return
	}

	aria2Cli = aria2srv.GetAria2()
	movie.NewMovieService(common.MovieSrv, client.DefaultClient)
	inited = true
}
