package loader

import (
	"fmt"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/hdt3213/godis/config"
	"github.com/hdt3213/godis/lib/utils"
	redisServer "github.com/hdt3213/godis/redis/server"
	"github.com/hdt3213/godis/tcp"
)

func startMemDB() {
	go startMySQL()

	go startRedis()

	time.Sleep(time.Second * 1)
}

var (
	dbName    = "test"
	tableName = "test"
	address   = "127.0.0.1"
	port      = 3306
)

func startMySQL() {
	ctx := sql.NewEmptyContext()

	db := memory.NewDatabase(dbName)
	db.EnablePrimaryKeyIndexes()
	table := memory.NewTable(tableName, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "id", Type: types.Int64, Nullable: false, Source: tableName, PrimaryKey: true},
		{Name: "goods_id", Type: types.TinyText, Nullable: false, Source: tableName},
		{Name: "name", Type: types.TinyText, Nullable: false, Source: tableName},
	}), db.GetForeignKeyCollection())
	db.AddTable(tableName, table)

	err := table.Insert(ctx, sql.NewRow(int64(1), "1", "golang"))
	if err != nil {
		panic(err)
	}

	engine := sqle.NewDefault(
		memory.NewDBProvider(
			db,
		))

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", address, port),
	}
	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}
	if err = s.Start(); err != nil {
		panic(err)
	}
}

func startRedis() {
	config.Properties = &config.ServerProperties{
		Bind:           "0.0.0.0",
		Port:           6379,
		AppendOnly:     false,
		AppendFilename: "",
		MaxClients:     1000,
		RunID:          utils.RandString(40),
	}
	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, redisServer.MakeHandler())
	if err != nil {
		panic(err)
	}
}
