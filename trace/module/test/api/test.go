package api

import (
	"time"

	"github.com/air-go/rpc/server/http/response"
	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/http"

	"github.com/air-go/rpc-example/trace/rpc/test"
)

func RPC(c *gin.Context) {
	time.Sleep(time.Millisecond * 30)
	ret, err := test.RPC(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.ErrnoServer, response.WrapToast(err.Error()))
		return
	}

	response.ResponseJSON(c, response.ErrnoSuccess, ret)
}

type RPC1Request struct {
	A string `json:"a"`
}

func RPC1(c *gin.Context) {
	time.Sleep(time.Millisecond * 99)
	var req RPC1Request
	if err := http.ParseAndValidateBody(c.Request, &req); err != nil {
		response.ResponseJSON(c, response.ErrnoParams, response.WrapToast(err.Error()))
		return
	}

	ret, err := test.RPC1(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.ErrnoServer, response.WrapToast(err.Error()))
		return
	}

	response.ResponseJSON(c, response.ErrnoSuccess, ret)
}

func Panic(c *gin.Context) {
	panic("test err")
}
