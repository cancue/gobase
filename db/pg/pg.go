package pg

import (
	"fmt"

	"github.com/cancue/gobase/config"
	"github.com/go-pg/pg/v9"
)

var db *pg.DB

// Set is exported.
func Set(conf *config.Config) {
	if yaml, ok := conf.YAML["db"].(map[string]interface{}); ok {
		connection := yaml["connection"].(map[string]interface{})

		db = pg.Connect(&pg.Options{
			Database: connection["database"].(string),
			Addr:     fmt.Sprintf("%s:%d", connection["host"], connection["port"]),
			User:     connection["user"].(string),
			Password: connection["password"].(string),
		})
	}
}

// Get returns go-pg DB which is safe for concurrent use by multiple goroutines and maintains its own connection pool.
func Get() *pg.DB {
	return db
}
