package provider

import (
	"go-framework/internal/rpc/helloworld"
	"go-framework/internal/server"
	"go-framework/proto/pb"
	"google.golang.org/grpc"
)

// RpcRegister 注册rpc服务
func RpcRegister(svc *server.SvcContext, rpcService grpc.ServiceRegistrar) {
	pb.RegisterSayWhatServer(rpcService, helloworld.NewSayWhatHandler(svc))
	//protobuf.RegisterSayWhatServer(rpcService, helloworld.NewSayWhatHandler())
	//protobuf.RegisterSayWhatServer(rpcService, new(helloworld2.SayWhatService))
}
