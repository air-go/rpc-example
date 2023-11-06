package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/air-go/rpc/bootstrap"
	client "github.com/air-go/rpc/client/http"
	"github.com/air-go/rpc/client/http/transport"
	"github.com/air-go/rpc/library/logger"
	"github.com/air-go/rpc/library/logger/zap"
	"github.com/air-go/rpc/library/servicer/load"
	"github.com/air-go/rpc/library/servicer/service"
	httpServer "github.com/air-go/rpc/server/http"
	logMiddleware "github.com/air-go/rpc/server/http/middleware/log"
	"github.com/air-go/rpc/server/http/response"
	"github.com/gin-gonic/gin"
	jsonCodec "github.com/why444216978/codec/json"
)

const (
	uri         = "/test"
	serviceName = "test"
	port        = 8777
)

func RegisterRouter(server *gin.Engine) {
	server.GET("/test", func(c *gin.Context) {
		response.ResponseJSON(c, http.StatusOK, nil, nil)
	})
}

func main() {
	go func() {
		srv := httpServer.New(fmt.Sprintf(":%d", port), RegisterRouter,
			httpServer.WithDebug(true),
			httpServer.WithMiddleware(logMiddleware.LoggerMiddleware(zap.StdLogger)),
		)
		if err := bootstrap.NewApp(srv).Start(); err != nil {
			log.Println(err)
			return
		}
	}()

	call()
}

type Response struct {
	Code    int         `json:"code"`
	Toast   string      `json:"toast"`
	Data    interface{} `json:"data"`
	Errmsg  string      `json:"errmsg"`
	TraceID string      `json:"trace_id"`
}

func call() {
	var err error
	if err = load.LoadService(&service.Config{
		ServiceName: serviceName,
		Type:        2,
		Host:        "127.0.0.1",
		Port:        port,
		Selector:    "wr",
	}); err != nil {
		log.Println(err)
		return
	}

	cli := transport.New(
		transport.WithLogger(zap.StdLogger))
	if err != nil {
		return
	}

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		resp := &client.DataResponse{
			Body:  new(map[string]interface{}),
			Codec: jsonCodec.JSONCodec{},
		}

		h := http.Header{}
		h.Set("Content-Type", "application/json")
		q := url.Values{}
		q.Add("test", "test")
		if err = cli.Send(logger.InitFieldsContainer(context.TODO()), &client.DefaultRequest{
			ServiceName: "test",
			Path:        uri,
			Method:      http.MethodGet,
			Codec:       jsonCodec.JSONCodec{},
			Header:      h,
			Query:       q,
		}, resp); err != nil {
			log.Fatal(err)
		}
		b, _ := json.Marshal(resp.Body)
		fmt.Println(string(b))
	}
}
