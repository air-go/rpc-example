package api

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/air-go/go-air-example/trace/response"
	"github.com/air-go/go-air-example/trace/rpc/test"
	httpResponse "github.com/air-go/rpc/server/http/response"
	"github.com/why444216978/go-util/http"
)

func RPC(c *gin.Context) {
	time.Sleep(time.Millisecond * 30)
	ret, err := test.RPC(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, httpResponse.WrapToast(err, err.Error()))
		return
	}

	response.ResponseJSON(c, response.CodeSuccess, ret, nil)
}

type RPC1Request struct {
	A string `json:"a"`
}

func RPC1(c *gin.Context) {
	time.Sleep(time.Millisecond * 99)
	var req RPC1Request
	if err := http.ParseAndValidateBody(c.Request, &req); err != nil {
		response.ResponseJSON(c, response.CodeParams, nil, httpResponse.WrapToast(err, err.Error()))
		return
	}

	ret, err := test.RPC1(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, httpResponse.WrapToast(err, err.Error()))
		return
	}

	response.ResponseJSON(c, response.CodeSuccess, ret, nil)
}

func Panic(c *gin.Context) {
	panic("test err")
}
