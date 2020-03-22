package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDbString(t *testing.T) {
	driver := "foo"
	yaml := map[string]interface{}{
		"driver": driver,
		"connection": map[string]interface{}{
			"database": "bar",
			"host":     "baz",
			"port":     1234,
			"user":     "qux",
			"password": "quux",
		},
	}
	expect := "port=1234 host=baz user=qux password=quux dbname=bar sslmode=disable"

	dbDriver, dbString := getDbString(yaml)

	assert.Equal(t, driver, dbDriver)
	assert.Equal(t, expect, dbString)
}
