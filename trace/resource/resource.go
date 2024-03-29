package resource

import (
	httpClient "github.com/air-go/rpc/client/http"
	"github.com/air-go/rpc/library/cache"
	"github.com/air-go/rpc/library/etcd"
	"github.com/air-go/rpc/library/lock"
	"github.com/air-go/rpc/library/logger"
	"github.com/air-go/rpc/library/orm"
	"github.com/air-go/rpc/library/queue"
	"github.com/air-go/rpc/library/redis"
	"github.com/air-go/rpc/library/registry"
)

var (
	TestDB        *orm.Orm
	RedisDefault  *redis.RedisClient
	Etcd          *etcd.Etcd
	ClientHTTP    httpClient.Client
	ServiceLogger logger.Logger
	RedisLock     lock.Locker
	RedisCache    cache.Cacher
	Registrar     registry.Registrar
	RabbitMQ      queue.Queue
)
