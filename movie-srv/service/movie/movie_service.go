package movie

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jlb0906/micro-movie/movie-srv/model"
	db2 "github.com/jlb0906/micro-movie/plugins/db"
	"sync"
)

var (
	m  sync.RWMutex
	s  *service
	db *gorm.DB
)

type Service interface {
	AddMovie(title, url, status string) (mid string, err error)
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

func (s *service) AddMovie(title, url, status string) (mid string, err error) {
	movie := model.Movie{
		Id:     uuid.New().String(),
		Title:  title,
		Url:    url,
		Status: status,
	}
	db.Create(&movie)
	b := db.NewRecord(movie)
	if b {
		return "", fmt.Errorf("插入错误")
	}
	return movie.Id, nil
}
