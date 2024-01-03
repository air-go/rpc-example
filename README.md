# go-air-example

# HTTP 服务测试
```
$cd http
#go run main.go 
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /test                     --> main.RegisterRouter.func1 (2 handlers)
{"level":"INFO","time":"2022-06-24 00:50:24","file":"log/log.go:85","func":"github.com/air-go/rpc/server/http/middleware/log.LoggerMiddleware.func1.2","msg":"request info","app_name":"","module":"default","service_name":"default","api":"/test","log_id":"1540014418616504320","header":{"Accept-Encoding":["gzip"],"Content-Length":["4"],"Log-Id":["1540014418616504320"],"Timeout-Millisecond":["0"],"User-Agent":["Go-http-client/1.1"]},"method":"GET","server_ip":"192.168.1.103","client_port":0,"server_port":0,"code":200,"cost":1,"trace_id":"","request":"R0VUIC90ZXN0IEhUVFAvMS4xDQpIb3N0OiAxMjcuMC4wLjE6ODc3Nw0KQWNjZXB0LUVuY29kaW5nOiBnemlwDQpDb250ZW50LUxlbmd0aDogNA0KTG9nLUlkOiAxNTQwMDE0NDE4NjE2NTA0MzIwDQpUaW1lb3V0LU1pbGxpc2Vjb25kOiAwDQpVc2VyLUFnZW50OiBHby1odHRwLWNsaWVudC8xLjENCg0KbnVsbA==","response":{"code":200,"data":{},"errmsg":"toast","toast":"toast","trace_id":""},"client_ip":"127.0.0.1"}
......
```
<br>

# gRPC 服务测试
```
$cd grpc
#go run main.go 
message:"why world"
message:"why world"
......

$curl localhost:8777/v1/example/echo
{"message":" world"}
```
<br>

# 完整业务开发测试
## 推荐代码结构
```
- rpc-example
    - trace //用户应用目录
        - conf //服务配置文件目录
        - dev
        - liantiao
        - online
        - qa
        - loader //资源加载
        - resource
        - resource.go //全局资源
        - response
        - response.go //http响应
        - router
            - router.go //路由定义和中间件注册
        - rpc //三方rpc调用封装
            - test //test服务
        - module //各模块核心实现，按照业务边界划分目录
        - module1 //模块1
            - api //对外暴露api
            - job //离线任务
            - responsitory //存储层
            - service //核心业务代码
        - module1 //模块2
            - api //对外暴露api
            - job //离线任务
            - responsitory //存储层
            - service //核心业务代码
        - main.go //app入口文件
```
<br>

## 准备工作
1. loader 下的资源初始化和服务注册是留给开发者自己扩展的，可自行调整资源加载。
2. 查看 conf/xxx 目录的各个配置文件，改成符合自己的
3. log 配置中的目录确保本地存在且有写入权限
4. 创建 test 数据库表并随意插入数据
```
CREATE TABLE `test` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `goods_id` bigint(20) unsigned NOT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_goods` (`goods_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin 
```
<br>

## 接口列表
检测接口：http://localhost:8777/ping 
<br>
panic 接口：http://localhost:8777/test/panic
<br>
分布式链路追踪：http://localhost:8777/ping/rpc
<br>
db 和 redis测试接口：http://localhost:8777/test/conn （依赖 mysql 和 redis）
<br>
分布式链路追踪+资源依赖：http://localhost:8777/test/rpc (依赖 mysql 和 redis）
<br>
<img src="https://github.com/why444216978/images/blob/master/jaeger.png" width="800" height="300" alt="jaeger"/>
<br>

## 运行
```
$cd trace
$go run main.go -env dev
2022/06/24 00:42:45 Actual pid is 82122
2022/06/24 00:42:45 start http, port 8777
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /debug/pprof/             --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/cmdline      --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/profile      --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] POST   /debug/pprof/symbol       --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/symbol       --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/trace        --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/allocs       --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/block        --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/goroutine    --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/heap         --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/mutex        --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /debug/pprof/threadcreate --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /ping                     --> github.com/air-go/rpc-example/trace/module/ping/api.Ping (4 handlers)
[GIN-debug] GET    /ping/rpc                 --> github.com/air-go/rpc-example/trace/module/ping/api.RPC (4 handlers)
[GIN-debug] POST   /test/rpc                 --> github.com/air-go/rpc-example/trace/module/test/api.RPC (4 handlers)
[GIN-debug] POST   /test/rpc1                --> github.com/air-go/rpc-example/trace/module/test/api.RPC1 (4 handlers)
[GIN-debug] POST   /test/panic               --> github.com/air-go/rpc-example/trace/module/test/api.Panic (4 handlers)
[GIN-debug] POST   /test/conn                --> github.com/air-go/rpc-example/trace/module/goods/api.Do (4 handlers)
watching prefix:rpc-example now...
service rpc-example  put key: rpc-example.192.168.1.104.8777 val: {"Host":"192.168.1.104","Port":8777}


