package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/air-go/rpc/bootstrap"
	"github.com/air-go/rpc/library/app"
	httpPrometheus "github.com/air-go/rpc/library/prometheus/http"
	httpServer "github.com/air-go/rpc/server/http"
	logMiddleware "github.com/air-go/rpc/server/http/middleware/log"
	panicMiddleware "github.com/air-go/rpc/server/http/middleware/panic"
	responseMiddleware "github.com/air-go/rpc/server/http/middleware/response"
	timeoutMiddleware "github.com/air-go/rpc/server/http/middleware/timeout"
	traceMiddleware "github.com/air-go/rpc/server/http/middleware/trace"

	"github.com/air-go/rpc-example/trace/loader"
	"github.com/air-go/rpc-example/trace/resource"
	"github.com/air-go/rpc-example/trace/router"
)

var env = flag.String("env", "dev", "config path")

func main() {
	flag.Parse()

	if err := bootstrap.Init("conf/"+*env, loader.Load); err != nil {
		log.Printf("bootstrap.Init err %s", err.Error())
		return
	}

	srv := httpServer.New(fmt.Sprintf(":%d", app.Port()), router.RegisterRouter,
		httpServer.WithReadTimeout(app.ReadTimeout()),
		httpServer.WithWriteTimeout(app.WriteTimeout()),
		httpServer.WithMiddleware(
			panicMiddleware.PanicMiddleware(resource.ServiceLogger),
			timeoutMiddleware.TimeoutMiddleware(app.ContextTimeout()),
			// traceMiddleware.OpentracingMiddleware(),
			traceMiddleware.OpentelemetryMiddleware(),
			logMiddleware.LoggerMiddleware(resource.ServiceLogger),
			httpPrometheus.HTTPMetricsMiddleware(),
			responseMiddleware.ResponseMiddleware(),
		),
		httpServer.WithPprof(app.Pprof()),
		httpServer.WithDebug(app.Debug()),
		httpServer.WithMetrics("/metrics"),
	)

	if err := bootstrap.NewApp(srv, bootstrap.WithRegistrar(resource.Registrar)).Start(); err != nil {
		log.Println(err)
	}
}
