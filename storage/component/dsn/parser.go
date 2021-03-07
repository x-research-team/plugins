package dsn

import (
	"fmt"
	"path/filepath"
	"sync"

	"entgo.io/ent/dialect"

	"github.com/x-research-team/bus"
	"github.com/x-research-team/utils/file"
)

var m sync.Mutex

type TDataBaseConfig struct {
	Name     string `json:"name"`
	Dialect  string `json:"dialect"`
	Database string `json:"database"`
}

func (c TDataBaseConfig) GetDialect() string {
	return c.Dialect
}

func (c TDataBaseConfig) GetDSN() string {
	return ""
}

type IDataBaseConfig interface {
	GetDialect() string
	GetDSN() string
}

type TMySQLConfig struct {
	TDataBaseConfig
}

type TSQLiteConfig struct {
	TDataBaseConfig
}

func (c TSQLiteConfig) GetDSN() string {
	return fmt.Sprintf("file:%s?cache=shared", c.TDataBaseConfig.Database)
}

type TPostgresConfig struct {
	TDataBaseConfig
}

func Parse() map[string]IDataBaseConfig {
	v := make(map[string]IDataBaseConfig)
	dbconfigs, err := filepath.Glob("**/*.dbconfig")
	if err != nil {
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(len(dbconfigs))
	for _, dbconfig := range dbconfigs {
		go func(wg *sync.WaitGroup, dbconfig string) {
			defer wg.Done()
			c := new(TDataBaseConfig)
			if err := file.Read(".", dbconfig, c); err != nil {
				bus.Error <- err
				return
			}
			if c.Name == "" {
				bus.Error <- fmt.Errorf("name of connection can not be empty")
				return
			}
			var s IDataBaseConfig
			switch c.GetDialect() {
			case dialect.MySQL:
				s = new(TMySQLConfig)
			case dialect.SQLite:
				s = new(TSQLiteConfig)
			case dialect.Postgres:
				s = new(TPostgresConfig)
			}
			if err := file.Read(".", dbconfig, s); err != nil {
				bus.Error <- err
				return
			}
			m.Lock()
			v[c.Name] = s
			m.Unlock()
		}(&wg, dbconfig)
	}
	wg.Wait()
	return v
}
