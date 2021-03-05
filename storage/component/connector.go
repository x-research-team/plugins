package component

import (
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/entc/integration/ent"
	"github.com/x-research-team/contract"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectTo(dsn string) contract.ComponentModule {
	return func(component contract.IComponent) {
		c := component.(*Component)
		drv, err := sql.Open(dialect.MySQL, dsn)
		if err != nil {
			c.fails = append(c.fails, err)
			return
		}
		db := drv.DB()
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
		c.client = ent.NewClient(ent.Driver(drv))
		component = c
	}
}
