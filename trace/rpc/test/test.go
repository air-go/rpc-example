package test

import (
	"context"
	"net/http"
	"net/url"

	jsonCodec "github.com/why444216978/codec/json"

	"github.com/air-go/go-air-example/trace/resource"
	httpClient "github.com/air-go/rpc/client/http"
	"github.com/air-go/rpc/library/logger"
)

const (
	serviceName = "rpc-example"
)

func RPC(ctx context.Context) (resp *httpClient.DataResponse, err error) {
	q := url.Values{}
	q.Add("logid", logger.ValueLogID(ctx))
	req := &httpClient.DefaultRequest{
		ServiceName: serviceName,
		Path:        "/test/rpc1",
		Query:       q,
		Method:      http.MethodPost,
		Header:      nil,
		Body:        map[string]interface{}{"rpc": "rpc"},
		Codec:       jsonCodec.JSONCodec{},
	}
	resp = &httpClient.DataResponse{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, req, resp); err != nil {
		return
	}

	return
}

func RPC1(ctx context.Context) (resp *httpClient.DataResponse, err error) {
	q := url.Values{}
	q.Add("logid", logger.ValueLogID(ctx))
	req := &httpClient.DefaultRequest{
		ServiceName: serviceName,
		Path:        "/test/conn",
		Query:       q,
		Method:      http.MethodPost,
		Header:      nil,
		Body:        map[string]interface{}{"rpc1": "rpc1"},
		Codec:       jsonCodec.JSONCodec{},
	}

	resp = &httpClient.DataResponse{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, req, resp); err != nil {
		return
	}

	return
}

func Ping(ctx context.Context) (resp *httpClient.DataResponse, err error) {
	q := url.Values{}
	q.Add("logid", logger.ValueLogID(ctx))
	req := &httpClient.DefaultRequest{
		ServiceName: serviceName,
		Path:        "/ping",
		Query:       q,
		Method:      http.MethodGet,
		Header:      nil,
		Body:        map[string]interface{}{"rpc1": "rpc1"},
		Codec:       jsonCodec.JSONCodec{},
	}

	resp = &httpClient.DataResponse{
		Body:  new(map[string]interface{}),
		Codec: jsonCodec.JSONCodec{},
	}

	if err = resource.ClientHTTP.Send(ctx, req, resp); err != nil {
		return
	}

	return
}
