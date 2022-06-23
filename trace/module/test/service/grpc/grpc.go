package grpc

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/air-go/go-air-example/trace/module/test/service/grpc/helloworld"
	serverGRPC "github.com/air-go/rpc/server/grpc"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: in.Name + " world"}, nil
}

func RegisterServer(s *grpc.Server) {
	pb.RegisterGreeterServer(s, &Server{})
}

func NewService() serverGRPC.Register {
	return serverGRPC.NewRegister(RegisterServer, pb.RegisterGreeterHandlerFromEndpoint)
}
