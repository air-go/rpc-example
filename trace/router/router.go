package router

import (
	"github.com/gin-gonic/gin"

	conn "github.com/air-go/rpc-example/trace/module/goods/api"
	ping "github.com/air-go/rpc-example/trace/module/ping/api"
	test "github.com/air-go/rpc-example/trace/module/test/api"
)

func RegisterRouter(server *gin.Engine) {
	pingGroup := server.Group("/ping")
	{
		pingGroup.GET("", ping.Ping)
		pingGroup.GET("/rpc", ping.RPC)
	}

	testGroup := server.Group("/test")
	{
		testGroup.POST("/rpc", test.RPC)
		testGroup.POST("/rpc1", test.RPC1)
		testGroup.POST("/panic", test.Panic)
		testGroup.POST("/conn", conn.Do)
	}
}
