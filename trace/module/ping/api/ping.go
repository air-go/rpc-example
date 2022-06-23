package api

import (
	"github.com/gin-gonic/gin"

	"github.com/air-go/go-air-example/trace/response"
	"github.com/air-go/go-air-example/trace/rpc/test"
	httpResponse "github.com/air-go/rpc/server/http/response"
)

func Ping(c *gin.Context) {
	response.ResponseJSON(c, response.CodeSuccess, nil, nil)
}

func RPC(c *gin.Context) {
	ret, err := test.Ping(c.Request.Context())
	if err != nil {
		response.ResponseJSON(c, response.CodeServer, ret, httpResponse.WrapToast(err, err.Error()))
		return
	}
	response.ResponseJSON(c, response.CodeSuccess, ret, nil)
}
