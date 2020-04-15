package handler

import (
	"context"
	"github.com/jlb0906/micro-movie/basic/common"
	proto "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	"github.com/jlb0906/micro-movie/movie-srv/service/movie"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	movieService movie.Service
)

type Orders struct {
}

// Init 初始化handler
func Init() {
	movieService, _ = movie.GetService()
}

type Movie struct{}

func (e *Movie) UpdateMovie(ctx context.Context, req *proto.UpdateReq, rsp *proto.UpdateRsp) error {
	movieService.UpdateMovie(req.Movie.Title, req.Movie.Url, req.Movie.Status)
	rsp.Msg = "succeed"
	return nil
}

func (e *Movie) AddMovie(ctx context.Context, req *proto.AddReq, resp *proto.AddRsp) error {
	_, err := movieService.AddMovie(req.Movie.Title, req.Movie.Url, req.Movie.Status)
	if err != nil {
		log.Error(err)
		resp.Msg = "failed"
		return errors.InternalServerError(common.MovieSrv, "内部错误")
	}

	resp.Msg = "succeed"
	return nil
}
