package aria2

import (
	"context"
	pb "github.com/jlb0906/micro-movie/movie-srv/proto/movie"
	aria2cfg "github.com/jlb0906/micro-movie/plugins/aria2"
	miniocfg "github.com/jlb0906/micro-movie/plugins/minio"
	"github.com/micro/go-micro/v2/logger"
	"github.com/minio/minio-go/v6"
	"github.com/zyxar/argo/rpc"
	"path/filepath"
	"strings"
	"time"
)

func createWorker(in chan []rpc.Event, s Scheduler) {
	go func() {
		for {
			s.WorkerReady(in)
			events := <-in
			worker(events)
		}
	}()
}

func worker(events []rpc.Event) {
	cli, conf := miniocfg.Get()

	for _, e := range events {
		stat, err := aria2Cli.TellStatus(e.Gid)
		if err != nil {
			logger.Error(err)
			continue
		}

		// Upload the file
		filePath := stat.Files[0].Path
		arr := strings.Split(filePath, "/")
		objectName := arr[len(arr)-1]
		filePath = filepath.Join(aria2cfg.Get().Prefix, filePath)

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

		movieCli.UpdateMovieByGid(context.TODO(), &pb.UpdateReq{
			Movie: &pb.MovieMsg{
				Gid:    e.Gid,
				Title:  objectName,
				Uri:    u.String(),
				Status: stat.Status,
			},
		})

		logger.Infof("Successfully uploaded %s of size %d", objectName, n)
	}
}
