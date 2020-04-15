package miniocfg

import (
	"github.com/jlb0906/micro-movie/basic"
	"github.com/jlb0906/micro-movie/basic/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/minio/minio-go/v6"
	"sync"
)

var (
	m      sync.RWMutex
	cfg    *MinioConf
	inited bool
	cli    *minio.Client
)

// Minio 配置
type MinioConf struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	UseSSL          bool   `json:"useSSL"`
	BucketName      string `json:"bucketName"`
	Location        string `json:"location"`
}

func init() {
	basic.Register(initMinio)
}

func initMinio() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Infof("[initMinio] 已经初始化过Minio...")
		return
	}

	log.Infof("[initMinio] 初始化Minio...")

	c := config.C()
	cfg = new(MinioConf)
	err := c.App("minio", cfg)
	if err != nil {
		log.Error("[initMinio] %s", err)
	}

	// Initialize minio client object.
	minioClient, err := minio.New(cfg.Endpoint, cfg.AccessKeyID, cfg.SecretAccessKey, cfg.UseSSL)
	if err != nil {
		log.Error(err)
	}

	err = minioClient.MakeBucket(cfg.BucketName, cfg.Location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(cfg.BucketName)
		if errBucketExists == nil && exists {
			log.Info("We already own %s", cfg.BucketName)
		} else {
			log.Error(err)
		}
	} else {
		log.Infof("Successfully created %s", cfg.BucketName)
	}

	cli = minioClient
	inited = true

	log.Infof("[initMinio] Minio，成功")
}

// Minio 获取Minio
func GetMinio() (*minio.Client, *MinioConf) {
	return cli, cfg
}
