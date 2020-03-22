package pg

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"

	"github.com/cancue/gobase/config"
)

// Goose is exported.
func Goose(conf *config.Config) {
	if yaml, ok := conf.YAML["db"].(map[string]interface{}); ok {

		dbDriver, dbString := getDbString(yaml)

		db, err := sql.Open(dbDriver, dbString)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		args := os.Args[1:]
		if err := goose.Run(args[0], db, "./db/migrations", args[1:]...); err != nil {
			panic(err)
		}
	}
}

func getDbString(yaml map[string]interface{}) (dbDriver string, dbString string) {
	connection := yaml["connection"].(map[string]interface{})

	dbDriver = yaml["driver"].(string)
	dbString = fmt.Sprintf(
		"port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
		connection["port"],
		connection["host"],
		connection["user"],
		connection["password"],
		connection["database"],
	)

	return
}
