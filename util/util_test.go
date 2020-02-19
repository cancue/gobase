package util

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	key := "gobase-test-key-dummy"
	value := "gobase-test-value-dummy"
	placeholder := "gobase-test-placeholder-dummy"
	os.Setenv(key, value)

	result := GetEnv(key, placeholder)

	assert.Equal(t, value, result)

	t.Run("placeholder", func(t *testing.T) {
		os.Unsetenv(key)

		result := GetEnv(key, placeholder)

		assert.Equal(t, placeholder, result)
	})
}

func TestReadYAMLFile(t *testing.T) {
	t.Run("should panic if file does not exist", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()

		ReadYAMLFile("")
	})

	t.Run("readYAML", func(t *testing.T) {
		file := strings.NewReader(`
a: Easy!
b:
  c: 2
  d:
  - 3
  - 4
`)
		expect := map[string]interface{}{
			"a": "Easy!",
			"b": map[string]interface{}{
				"c": 2,
				"d": []interface{}{
					3,
					4,
				},
			},
		}

		result := readYAML(file)

		assert.Equal(t, expect, result)
	})
}
