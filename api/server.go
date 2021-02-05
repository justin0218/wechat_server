package api

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"wechat_server/api/proto"
	"wechat_server/store"
)

type WechatSvr struct {
	conf store.Config
	redis store.Redis
}

func GrpcServer() {
	conf := new(store.Config)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", conf.Get().Etcd.Key))
	if err != nil {
		panic(err)
	}
	var opts []grpc.ServerOption
	svr := grpc.NewServer(opts...)
	proto.RegisterWechatServer(svr, &WechatSvr{})
	err = svr.Serve(lis)
	if err != nil {
		panic(err)
	}
}
