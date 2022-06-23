package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/air-go/go-air-example/trace/module/test/service/grpc/helloworld"
	client "github.com/air-go/rpc/client/grpc"
	"github.com/air-go/rpc/library/app"
)

func Start(ctx context.Context) (err error) {
	call()
	return
}

func call() {
	cc, err := client.Conn(context.Background(), fmt.Sprintf(":%d", app.Port()))
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
