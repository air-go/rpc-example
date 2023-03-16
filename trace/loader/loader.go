package loader

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/why444216978/go-util/assert"
	"github.com/why444216978/go-util/sys"

	"github.com/air-go/go-air-example/trace/resource"
	httpClient "github.com/air-go/rpc/client/http"
	"github.com/air-go/rpc/client/http/transport"
	"github.com/air-go/rpc/library/app"
	redisCache "github.com/air-go/rpc/library/cache/redis"
	"github.com/air-go/rpc/library/config"
	"github.com/air-go/rpc/library/etcd"
	redisLock "github.com/air-go/rpc/library/lock/redis"
	loggerGorm "github.com/air-go/rpc/library/logger/zap/gorm"
	loggerRedis "github.com/air-go/rpc/library/logger/zap/redis"
	loggerRPC "github.com/air-go/rpc/library/logger/zap/rpc"
	serviceLogger "github.com/air-go/rpc/library/logger/zap/service"
	"github.com/air-go/rpc/library/opentracing"
	"github.com/air-go/rpc/library/orm"
	"github.com/air-go/rpc/library/otel"
	otelJaeger "github.com/air-go/rpc/library/otel/exporters/jaeger"
	otelGorm "github.com/air-go/rpc/library/otel/gorm"
	otelRedis "github.com/air-go/rpc/library/otel/redis"
	"github.com/air-go/rpc/library/queue/rabbitmq"
	"github.com/air-go/rpc/library/redis"
	etcdRegistry "github.com/air-go/rpc/library/registry/etcd"
	"github.com/air-go/rpc/library/servicer/load"
	"github.com/air-go/rpc/server"
)

// Load load resource
func Load() (err error) {
	// TODO
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	if err = loadLogger(); err != nil {
		return
	}
	if err = loadClientHTTP(); err != nil {
		return
	}
	if err = loadMysql("test_mysql"); err != nil {
		return
	}
	if err = loadRedis("default_redis"); err != nil {
		return
	}
	if err = loadJaeger(); err != nil {
		return
	}
	if err = loadOpentelemetry(); err != nil {
		return
	}
	if err = loadLock(); err != nil {
		return
	}
	if err = loadCache(); err != nil {
		return
	}
	if err = loadRabbitMQ("default_rabbitmq"); err != nil {
		return
	}
	if err = loadEtcd(); err != nil {
		return
	}
	if err = loadRegistry(); err != nil {
		return
	}
	if err = load.LoadGlobPattern("services", "toml", resource.Etcd); err != nil {
		return
	}

	return
}

func loadLogger() (err error) {
	cfg := &serviceLogger.Config{}

	if err = config.ReadConfig("log/service", "toml", &cfg); err != nil {
		return
	}

	if resource.ServiceLogger, err = serviceLogger.NewServiceLogger(app.Name(), cfg); err != nil {
		return
	}

	server.RegisterCloseFunc(resource.ServiceLogger.Close())

	return
}

func loadMysql(db string) (err error) {
	cfg := &orm.Config{}
	logCfg := &loggerGorm.GormConfig{}

	if err = config.ReadConfig(db, "toml", cfg); err != nil {
		return
	}

	if err = config.ReadConfig("log/gorm", "toml", logCfg); err != nil {
		return
	}

	logCfg.ServiceName = cfg.ServiceName
	logger, err := loggerGorm.NewGorm(logCfg)
	if err != nil {
		return
	}
	server.RegisterCloseFunc(logger.Close())

	if resource.TestDB, err = orm.NewOrm(cfg,
		// orm.WithTrace(jaegerGorm.GormTrace),
		orm.WithTrace(otelGorm.NewOpentelemetryPlugin()),
		orm.WithLogger(logger),
	); err != nil {
		return
	}

	return
}

func loadRedis(db string) (err error) {
	cfg := &redis.Config{}
	logCfg := &loggerRedis.RedisConfig{}

	if err = config.ReadConfig(db, "toml", cfg); err != nil {
		return
	}
	if err = config.ReadConfig("log/redis", "toml", &logCfg); err != nil {
		return
	}

	logCfg.ServiceName = cfg.ServiceName
	logCfg.Host = cfg.Host
	logCfg.Port = cfg.Port

	logger, err := loggerRedis.NewRedisLogger(logCfg)
	if err != nil {
		return
	}
	server.RegisterCloseFunc(logger.Close())

	rc, err := redis.NewRedisClient(cfg)
	if err != nil {
		return
	}
	// rc.AddHook(jaegerRedis.NewJaegerHook())
	rc.AddHook(otelRedis.NewOpentelemetryHook())
	rc.AddHook(logger)
	resource.RedisDefault = rc

	return
}

func loadRabbitMQ(service string) (err error) {
	cfg := &rabbitmq.Config{}
	if err = config.ReadConfig(service, "toml", cfg); err != nil {
		return
	}

	if resource.RabbitMQ, err = rabbitmq.New(cfg); err != nil {
		return
	}

	return
}

func loadLock() (err error) {
	resource.RedisLock, err = redisLock.New(resource.RedisDefault)
	return
}

func loadCache() (err error) {
	resource.RedisCache, err = redisCache.New(resource.RedisDefault, resource.RedisLock)
	return
}

func loadJaeger() (err error) {
	cfg := &opentracing.Config{}

	if err = config.ReadConfig("jaeger", "toml", cfg); err != nil {
		return
	}

	if _, _, err = opentracing.NewJaegerTracer(cfg, app.Name()); err != nil {
		return
	}

	return
}

func loadOpentelemetry() (err error) {
	cfg := &otelJaeger.JaegerConfig{}

	if err = config.ReadConfig("jaeger", "toml", cfg); err != nil {
		return
	}

	a, err := otelJaeger.NewJaeger(cfg)
	if err != nil {
		return
	}
	if err = otel.NewTracer(app.Name(), a.Exporter, otel.WithSampler(a.Sampler)); err != nil {
		return
	}
	return
}

func loadEtcd() (err error) {
	cfg := &etcd.Config{}

	if err = config.ReadConfig("etcd", "toml", cfg); err != nil {
		return
	}

	if resource.Etcd, err = etcd.NewClient(
		strings.Split(cfg.Endpoints, ";"),
		etcd.WithDialTimeout(cfg.DialTimeout),
	); err != nil {
		return
	}

	return
}

func loadRegistry() (err error) {
	var localIP string

	if localIP, err = sys.LocalIP(); err != nil {
		return
	}

	if assert.IsNil(resource.Etcd) {
		err = errors.New("resource.Etcd is nil")
		return
	}

	if resource.Registrar, err = etcdRegistry.NewRegistry(resource.Etcd.Client, app.RegistryName(), localIP, app.Port()); err != nil {
		return
	}

	if err = server.RegisterCloseFunc(resource.Registrar.DeRegister); err != nil {
		return
	}

	return
}

func loadClientHTTP() (err error) {
	cfg := &loggerRPC.RPCConfig{}
	if err = config.ReadConfig("log/rpc", "toml", cfg); err != nil {
		return
	}

	logger, err := loggerRPC.NewRPCLogger(cfg)
	if err != nil {
		return
	}
	server.RegisterCloseFunc(logger.Close())

	resource.ClientHTTP = transport.New(
		transport.WithLogger(logger),
		// transport.WithBeforePlugins(&httpClient.OpentracingBeforePlugin{}),
		transport.WithBeforePlugins(&httpClient.OpentelemetryOpentracingBeforePlugin{}))
	if err != nil {
		return
	}

	return
}
