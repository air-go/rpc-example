package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/air-go/go-air-example/trace/loader"
	jobGRPC "github.com/air-go/go-air-example/trace/module/test/job/grpc"
	"github.com/air-go/go-air-example/trace/resource"
	"github.com/air-go/go-air-example/trace/router"
	"github.com/air-go/rpc/bootstrap"
	"github.com/air-go/rpc/library/app"
	jobLib "github.com/air-go/rpc/library/job"
	httpServer "github.com/air-go/rpc/server/http"
	logMiddleware "github.com/air-go/rpc/server/http/middleware/log"
	panicMiddleware "github.com/air-go/rpc/server/http/middleware/panic"
	timeoutMiddleware "github.com/air-go/rpc/server/http/middleware/timeout"
	traceMiddleware "github.com/air-go/rpc/server/http/middleware/trace"
)

var (
	env = flag.String("env", "dev", "config path")
	job = flag.String("job", "", "is job")
)

func main() {
	flag.Parse()

	var err error

	if err = bootstrap.Init("conf/"+*env, loader.Load); err != nil {
		log.Printf("bootstrap.Init err %s", err.Error())
		return
	}

	if *job != "" {
		jobLib.Handlers = map[string]jobLib.HandleFunc{
			"grpc-test": jobGRPC.Start,
		}
		jobLib.Handle(*job, resource.ServiceLogger)
		return
	}

	srv := httpServer.New(fmt.Sprintf(":%d", app.Port()), router.RegisterRouter,
		httpServer.WithReadTimeout(app.ReadTimeout()),
		httpServer.WithWriteTimeout(app.WriteTimeout()),
		httpServer.WithMiddleware(
			panicMiddleware.ThrowPanic(resource.ServiceLogger),
			timeoutMiddleware.TimeoutMiddleware(app.ContextTimeout()),
			// traceMiddleware.OpentracingMiddleware(),
			traceMiddleware.OpentelemetryMiddleware(),
			logMiddleware.LoggerMiddleware(resource.ServiceLogger),
		),
		httpServer.WithPprof(app.Pprof()),
		httpServer.WithDebug(app.Debug()),
	)

	if err := bootstrap.NewApp(srv, bootstrap.WithRegistrar(resource.Registrar)).Start(); err != nil {
		log.Println(err)
	}
}
