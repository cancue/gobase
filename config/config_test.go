package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	yamlFiles := map[string][]byte{
		"local": []byte("name: gobase\nserver:\n domain: localhost\n port: 3000\n allow-origins:\n - http://localhost:3000\n - https://example.com"),
	}
	SetYAMLs(yamlFiles)
	config := Get()

	t.Run("properties", func(t *testing.T) {
		assert.Equal(t, "local", config.Stage)
		assert.Equal(t, "gobase", config.Name)
		assert.Equal(t, "localhost", config.Domain)
		assert.Equal(t, 3000, config.Port)

		allowOrigins := config.AllowOrigins
		assert.Equal(t, 2, len(allowOrigins))
		assert.Equal(t, "http://localhost:3000", allowOrigins[0])
		assert.Equal(t, "https://example.com", allowOrigins[1])
	})
}

func TestGetEnv(t *testing.T) {
	key := "foo"
	value := "bar"
	placeholder := "baz"

	t.Run("envkey", func(t *testing.T) {
		os.Setenv(key, value)

		result := getEnv(key, placeholder)

		assert.Equal(t, value, result)
	})

	t.Run("placeholder", func(t *testing.T) {
		os.Unsetenv(key)

		result := getEnv(key, placeholder)

		assert.Equal(t, placeholder, result)
	})
}
