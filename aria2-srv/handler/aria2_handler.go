package handler

import (
	"context"
	"fmt"
	aria2 "github.com/jlb0906/micro-movie/aria2-srv/proto/aria2"
	aria22 "github.com/jlb0906/micro-movie/aria2-srv/service/aria2"
	"github.com/jlb0906/micro-movie/basic/common"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/zyxar/argo/rpc"
)

var aria2Cli rpc.Client

type Aria2 struct{}

func (e *Aria2) AddURI(ctx context.Context, req *aria2.AddURIReq, rsp *aria2.AddURIRsp) error {
	logger.Info("Received Aria2.AddURI request")

	gid, err := aria2Cli.AddURI(req.Uri)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &aria2.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("添加了下载任务：%v", gid)

	rsp.Gid = gid
	return nil
}

func (e *Aria2) Remove(ctx context.Context, req *aria2.RemoveReq, rsp *aria2.RemoveRsp) error {
	logger.Info("Received Aria2.Remove request")

	gid, err := aria2Cli.Remove(req.Gid)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &aria2.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("删除了下载任务：%v", gid)

	rsp.Gid = "gid"
	return nil
}

func (e *Aria2) Pause(ctx context.Context, req *aria2.PauseReq, rsp *aria2.PauseRsp) error {
	logger.Info("Received Aria2.Pause request")

	gid, err := aria2Cli.Pause(req.Gid)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &aria2.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("暂停了下载任务：%v", gid)

	rsp.Gid = "gid"
	return nil
}

func (e *Aria2) TellStatus(ctx context.Context, req *aria2.TellStatusReq, rsp *aria2.TellStatusRsp) error {
	logger.Info("Received Aria2.TellStatus request")

	info, err := aria2Cli.TellStatus(req.Gid, req.Keys...)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf("错误的请求 %v", err)
		rsp.Err = &aria2.Error{
			Code:   400,
			Detail: msg,
		}
		return errors.BadRequest(common.Aria2Srv, msg)
	}
	logger.Infof("下载任务的状态：%v", info)

	rsp.Info = &aria2.StatusInfo{
		Gid:             info.Gid,
		Status:          info.Status,
		TotalLength:     info.TotalLength,
		CompletedLength: info.CompletedLength,
		UploadLength:    info.UploadLength,
		BitField:        info.BitField,
		DownloadSpeed:   info.DownloadSpeed,
		UploadSpeed:     info.UploadSpeed,
		InfoHash:        info.InfoHash,
		NumSeeders:      info.NumSeeders,
		Connections:     info.Connections,
		ErrorCode:       info.ErrorCode,
		ErrorMessage:    info.ErrorMessage,
		FollowedBy:      info.FollowedBy,
		BelongsTo:       info.BelongsTo,
		Dir:             info.Dir,
		Files:           toProto(info.Files),
	}
	return nil
}

func toProto(files []rpc.FileInfo) []*aria2.FileInfo {
	infos := make([]*aria2.FileInfo, 0)
	for _, f := range files {
		uriInfos := make([]*aria2.URIInfo, 0)
		for _, u := range f.URIs {
			uriInfos = append(uriInfos, &aria2.URIInfo{
				URI:    u.URI,
				Status: u.Status,
			})
		}
		infos = append(infos, &aria2.FileInfo{
			Index:           f.Index,
			Path:            f.Path,
			Length:          f.Length,
			CompletedLength: f.CompletedLength,
			Selected:        f.Selected,
			URIs:            uriInfos,
		})
	}
	return infos
}

func Init() {
	aria2Cli = aria22.GetAria2(context.TODO())
}
