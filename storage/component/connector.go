package component

import (
	"storage/component/dsn"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/x-research-team/contract"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectTo(dsn map[string]dsn.IDataBaseConfig) contract.ComponentModule {
	return func(component contract.IComponent) {
		c := component.(*Component)
		for k, v := range dsn {
			drv, err := sql.Open(v.GetDialect(), v.GetDSN())
			if err != nil {
				c.fails = append(c.fails, err)
				return
			}
			db := drv.DB()
			db.SetMaxIdleConns(10)
			db.SetMaxOpenConns(100)
			db.SetConnMaxLifetime(time.Hour)
			if err = db.Ping(); err != nil {
				c.fails = append(c.fails, err)
				return
			}
			c.client[k] = db
		}
		component = c
	}
}
