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

func (e *Movie) AddMovie(ctx context.Context, req *proto.AddRequest, resp *proto.AddResponse) error {
	_, err := movieService.AddMovie(req.Title, req.Url, req.Status)
	if err != nil {
		log.Error(err)
		resp.Msg = "failed"
		return errors.InternalServerError(common.MovieSrv, "内部错误")
	}

	resp.Msg = "succeed"
	return nil
}
