package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/air-go/rpc/bootstrap"
	client "github.com/air-go/rpc/client/grpc"
	serverGRPC "github.com/air-go/rpc/server/grpc"
	serverH2C "github.com/air-go/rpc/server/grpc/h2c"
	"google.golang.org/grpc"

	pb "github.com/air-go/go-air-example/grpc/helloworld"
)

const endpoint = ":8777"

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

func main() {
	go func() {
		srv := serverH2C.NewH2C(endpoint, []serverGRPC.Register{NewService()})
		if err := bootstrap.NewApp(srv).Start(); err != nil {
			log.Println(err)
			return
		}
	}()

	call()
}

func call() {
	cc, err := client.Conn(context.Background(), endpoint)
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewGreeterClient(cc)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "why"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply)
	}
}
