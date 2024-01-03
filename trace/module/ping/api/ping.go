package api

import (
	"github.com/gin-gonic/gin"

	"github.com/air-go/rpc-example/trace/rpc/test"
	"github.com/air-go/rpc/server/http/response"
)

func Ping(c *gin.Context) {
	response.ResponseJSON(c, response.ErrnoSuccess)
}

func RPC(c *gin.Context) {
	ret, err := test.Ping(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.ErrnoServer, ret, response.WrapToast(err.Error()))
		return
	}
	response.ResponseJSON(c, response.ErrnoSuccess, ret)
}
