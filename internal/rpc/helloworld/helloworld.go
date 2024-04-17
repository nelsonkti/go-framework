package helloworld

import (
	"context"
	"fmt"
	"go-framework/internal/server"
	"go-framework/proto/pb"
	"go-framework/util/xerror"
)

type SayWhatHandler struct {
	svc                           *server.SvcContext
	pb.UnimplementedSayWhatServer // Embedding for forward compatibility
}

func NewSayWhatHandler(svc *server.SvcContext) *SayWhatHandler {
	return &SayWhatHandler{svc: svc}
}

func (s *SayWhatHandler) SayHello(ctx context.Context, request *pb.SayRequest) (*pb.SayResponse, error) {
	fmt.Println(request.Name)
	s.svc.Container.OrderService.Get()

	return &pb.SayResponse{Name: "fzy"}, xerror.BadRequest(3002, "操作问题")
}

func (s *SayWhatHandler) SayHello2(ctx context.Context, request *pb.SayRequest) (*pb.SayResponse, error) {
	fmt.Println(request.Name)
	//sw.svc.Container.OrderService.Get()
	return &pb.SayResponse{Name: "fzy2"}, nil
}
