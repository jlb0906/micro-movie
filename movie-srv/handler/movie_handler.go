package handler

import (
	"context"
	"fmt"
	aria2pb "github.com/jlb0906/micro-movie/aria2-srv/proto/aria2"
	"github.com/jlb0906/micro-movie/basic/common"
	pb "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	"github.com/jlb0906/micro-movie/movie-srv/service/movie"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"sync"
)

var (
	m            sync.RWMutex
	inited       bool
	movieService movie.Service
	aria2Srv     aria2pb.Aria2Service
)

type Orders struct {
}

// Init 初始化handler
func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		logger.Warn(fmt.Sprint("[Init] handler 已经初始化过"))
		return
	}

	movieService, _ = movie.GetService()
	aria2Srv = aria2pb.NewAria2Service(common.Aria2Srv, client.DefaultClient)
	inited = true
}

type Movie struct{}

func (e *Movie) SelectAll(ctx context.Context, _ *pb.Empty, rsp *pb.SelectRsp) error {
	rsp.Movies = movieService.SelectAll()
	return nil
}

func (e *Movie) UpdateMovie(ctx context.Context, req *pb.UpdateReq, rsp *pb.UpdateRsp) error {
	movieService.UpdateMovie(req.Movie)
	rsp.Msg = "succeed"
	return nil
}

func (e *Movie) UpdateMovieByGid(ctx context.Context, req *pb.UpdateReq, rsp *pb.UpdateRsp) error {
	movieService.UpdateMovieByGid(req.Movie)
	rsp.Msg = "succeed"
	return nil
}

func (e *Movie) AddMovie(ctx context.Context, req *pb.AddReq, resp *pb.AddRsp) error {
	_, err := movieService.AddMovie(req.Movie)
	if err != nil {
		logger.Error(err)
		resp.Msg = "failed"
		return errors.InternalServerError(common.MovieSrv, "内部错误")
	}

	aria2Srv.AddURI(ctx, &aria2pb.AddURIReq{Uri: req.Movie.Uri})
	resp.Msg = "succeed"
	return nil
}
