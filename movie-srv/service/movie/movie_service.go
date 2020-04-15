package movie

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jlb0906/micro-movie/movie-srv/model"
	proto "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	db2 "github.com/jlb0906/micro-movie/plugins/db"
	"sync"
)

var (
	m  sync.RWMutex
	s  *service
	db *gorm.DB
)

type Service interface {
	AddMovie(m *proto.Movie) (mid string, err error)
	UpdateMovie(m *proto.Movie)
	SelectAll() []*proto.Movie
}

// GetService 获取服务类
func GetService() (Service, error) {
	if s == nil {
		return nil, fmt.Errorf("[GetService] GetService 未初始化")
	}
	return s, nil
}

// Init 初始化库存服务层
func Init() {
	m.Lock()
	defer m.Unlock()

	if s != nil {
		return
	}
	s = &service{}
	db = db2.GetDB()

	if !db.HasTable(&model.Movie{}) {
		db.CreateTable(&model.Movie{})
	}
}

// service 服务
type service struct {
}

func (s *service) SelectAll() []*proto.Movie {
	var movies []*proto.Movie
	db.Find(&movies)
	return movies
}

func (s *service) UpdateMovie(m *proto.Movie) {
	db.Model(&model.Movie{Id: m.Id}).Updates(map[string]interface{}{"title": m.Title, "url": m.Url, "status": m.Status})
}

func (s *service) AddMovie(m *proto.Movie) (mid string, err error) {
	movie := model.Movie{
		Id:     m.Id,
		Title:  m.Title,
		Url:    m.Url,
		Status: m.Status,
	}
	db.Create(&movie)
	b := db.NewRecord(movie)
	if b {
		return "", fmt.Errorf("插入错误")
	}
	return movie.Id, nil
}
