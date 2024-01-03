package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/air-go/rpc/bootstrap"
	client "github.com/air-go/rpc/client/grpc"
	"github.com/air-go/rpc/library/selector/factory"
	"github.com/air-go/rpc/library/servicer/load"
	"github.com/air-go/rpc/library/servicer/service"
	serverGRPC "github.com/air-go/rpc/server/grpc"
	h2c "github.com/air-go/rpc/server/grpc/h2c"
	"google.golang.org/grpc"

	pb "github.com/air-go/rpc-example/grpc/helloworld"
)

const endpoint = ":8777"

const serviceName = "rpc-example"

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: in.Name + " world"}, nil
}

func RegisterServer(s *grpc.Server) {
	pb.RegisterGreeterServer(s, &Server{})
}

func main() {
	go func() {
		ctx := context.Background()
		h, err := serverGRPC.RegisterGateway(ctx, endpoint, []serverGRPC.RegisterMux{pb.RegisterGreeterHandlerFromEndpoint})
		if err != nil {
			panic(err)
		}

		srv := h2c.NewH2C(ctx, endpoint, []serverGRPC.RegisterGRPC{RegisterServer}, h2c.WithHTTPHandler(h))
		if err := bootstrap.NewApp(srv).Start(); err != nil {
			log.Println(err)
			return
		}
	}()

	cfg := &service.Config{
		ServiceName:  serviceName,
		RegistryName: "rpc-example-dev",
		Type:         2,
		Host:         "127.0.0.1",
		Port:         8777,
		Selector:     "wr",
	}
	if err := load.LoadService(cfg, service.WithSelector(factory.New(cfg.ServiceName, cfg.Selector))); err != nil {
		log.Println(err)
		return
	}

	// wait server start
	time.Sleep(time.Second * 3)

	call()
}

func call() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cc, err := client.Conn(ctx, serviceName)
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	client := pb.NewGreeterClient(cc)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "why"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply)
	}
}
