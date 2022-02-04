package data

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	_ "github.com/lib/pq"
	"github.com/menta2l/valhalla-layout/internal/conf"
	"github.com/menta2l/valhalla-layout/internal/data/ent"
	"github.com/menta2l/valhalla-layout/internal/data/ent/migrate"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewRegister, NewEntClient)

// Data .
type Data struct {
	db     *ent.Client
	log    *log.Helper
	secret string
}

// NewData .
func NewData(c *conf.Data, entClient *ent.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
func NewRegister(conf *conf.Registry) registry.Registrar {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	cli.Connect()
	return r
}
func NewEntClient(conf *conf.Data, logger log.Logger) *ent.Client {
	log := log.NewHelper(log.With(logger, "module", "data/ent"))
	db, err := sql.Open(dialect.Postgres, conf.Database.Source)
	if err != nil {
		log.Fatal(err)
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	drvWithContext := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		log.Infof("%v", i)
		// attach tracing here
		// Example output:
		// 2022/01/20 00:14:06 [Tx(4fe0fa85-1027-475b-93a8-d22c9cef287c).Query: query=SELECT COUNT(*) FROM `sqlite_master` WHERE `type` = ? AND `name` = ? args=[table ent_types]]
	})
	client := ent.NewClient(ent.Driver(drvWithContext))
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
