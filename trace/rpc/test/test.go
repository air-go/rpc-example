package test

import (
	"context"
	"fmt"
	"net/http"

	jsonCodec "github.com/why444216978/codec/json"

	"github.com/air-go/go-air-example/trace/resource"
	httpClient "github.com/air-go/rpc/client/http"
	"github.com/air-go/rpc/library/logger"
)

const (
	serviceName = "rpc-example"
)

func RPC(ctx context.Context) (resp *httpClient.Response, err error) {
	req := httpClient.Request{
		URI:    fmt.Sprintf("/test/rpc1?logid=%s", logger.ValueLogID(ctx)),
		Method: http.MethodPost,
		Header: nil,
		Body:   map[string]interface{}{"rpc": "rpc"},
		Codec:  jsonCodec.JSONCodec{},
	}
	resp = &httpClient.Response{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, serviceName, req, resp); err != nil {
		return
	}

	return
}

func RPC1(ctx context.Context) (resp *httpClient.Response, err error) {
	req := httpClient.Request{
		URI:    fmt.Sprintf("/test/conn?logid=%s", logger.ValueLogID(ctx)),
		Method: http.MethodPost,
		Header: nil,
		Body:   map[string]interface{}{"rpc1": "rpc1"},
		Codec:  jsonCodec.JSONCodec{},
	}

	resp = &httpClient.Response{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, serviceName, req, resp); err != nil {
		return
	}

	return
}

func Ping(ctx context.Context) (resp *httpClient.Response, err error) {
	req := httpClient.Request{
		URI:    fmt.Sprintf("/ping?logid=%s", logger.ValueLogID(ctx)),
		Method: http.MethodGet,
		Header: nil,
		Body:   map[string]interface{}{"rpc1": "rpc1"},
		Codec:  jsonCodec.JSONCodec{},
	}

	resp = &httpClient.Response{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, serviceName, req, resp); err != nil {
		return
	}

	return
}
