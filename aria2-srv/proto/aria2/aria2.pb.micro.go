// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: aria2-srv/proto/aria2/aria2.proto

package aria2

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Aria2 service

func NewAria2Endpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Aria2 service

type Aria2Service interface {
	AddURI(ctx context.Context, in *AddURIReq, opts ...client.CallOption) (*AddURIRsp, error)
	Remove(ctx context.Context, in *RemoveReq, opts ...client.CallOption) (*RemoveRsp, error)
	Pause(ctx context.Context, in *PauseReq, opts ...client.CallOption) (*PauseRsp, error)
	TellStatus(ctx context.Context, in *TellStatusReq, opts ...client.CallOption) (*TellStatusRsp, error)
}

type aria2Service struct {
	c    client.Client
	name string
}

func NewAria2Service(name string, c client.Client) Aria2Service {
	return &aria2Service{
		c:    c,
		name: name,
	}
}

func (c *aria2Service) AddURI(ctx context.Context, in *AddURIReq, opts ...client.CallOption) (*AddURIRsp, error) {
	req := c.c.NewRequest(c.name, "Aria2.AddURI", in)
	out := new(AddURIRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aria2Service) Remove(ctx context.Context, in *RemoveReq, opts ...client.CallOption) (*RemoveRsp, error) {
	req := c.c.NewRequest(c.name, "Aria2.Remove", in)
	out := new(RemoveRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aria2Service) Pause(ctx context.Context, in *PauseReq, opts ...client.CallOption) (*PauseRsp, error) {
	req := c.c.NewRequest(c.name, "Aria2.Pause", in)
	out := new(PauseRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aria2Service) TellStatus(ctx context.Context, in *TellStatusReq, opts ...client.CallOption) (*TellStatusRsp, error) {
	req := c.c.NewRequest(c.name, "Aria2.TellStatus", in)
	out := new(TellStatusRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Aria2 service

type Aria2Handler interface {
	AddURI(context.Context, *AddURIReq, *AddURIRsp) error
	Remove(context.Context, *RemoveReq, *RemoveRsp) error
	Pause(context.Context, *PauseReq, *PauseRsp) error
	TellStatus(context.Context, *TellStatusReq, *TellStatusRsp) error
}

func RegisterAria2Handler(s server.Server, hdlr Aria2Handler, opts ...server.HandlerOption) error {
	type aria2 interface {
		AddURI(ctx context.Context, in *AddURIReq, out *AddURIRsp) error
		Remove(ctx context.Context, in *RemoveReq, out *RemoveRsp) error
		Pause(ctx context.Context, in *PauseReq, out *PauseRsp) error
		TellStatus(ctx context.Context, in *TellStatusReq, out *TellStatusRsp) error
	}
	type Aria2 struct {
		aria2
	}
	h := &aria2Handler{hdlr}
	return s.Handle(s.NewHandler(&Aria2{h}, opts...))
}

type aria2Handler struct {
	Aria2Handler
}

func (h *aria2Handler) AddURI(ctx context.Context, in *AddURIReq, out *AddURIRsp) error {
	return h.Aria2Handler.AddURI(ctx, in, out)
}

func (h *aria2Handler) Remove(ctx context.Context, in *RemoveReq, out *RemoveRsp) error {
	return h.Aria2Handler.Remove(ctx, in, out)
}

func (h *aria2Handler) Pause(ctx context.Context, in *PauseReq, out *PauseRsp) error {
	return h.Aria2Handler.Pause(ctx, in, out)
}

func (h *aria2Handler) TellStatus(ctx context.Context, in *TellStatusReq, out *TellStatusRsp) error {
	return h.Aria2Handler.TellStatus(ctx, in, out)
}
