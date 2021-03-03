package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config .
type Config struct {
	Stage        string
	Name         string
	Domain       string
	Port         int
	AllowOrigins []string
	YAML         map[string]interface{}
}

// Get retrieves config.
func Get() *Config {
	return config
}

// Set set config to use globally.
func Set(conf *Config) {
	if config != nil {
		panic("config already set")
	}

	action.Do(func() {
		config = conf
	})
}

// SetYAMLs set config from stage-yaml mapped argument.
func SetYAMLs(yamlFiles map[string][]byte) {
	stage := getEnv("STAGE", "local")
	yaml := unmarshalYAML(yamlFiles[stage])
	name := yaml["name"].(string)
	server := yaml["server"].(map[string]interface{})
	port := server["port"].(int)
	domain := server["domain"].(string)

	var allowOrigins []string
	for _, v := range server["allow-origins"].([]interface{}) {
		allowOrigins = append(allowOrigins, v.(string))
	}

	Set(&Config{
		Stage:        stage,
		Name:         name,
		Domain:       domain,
		Port:         port,
		AllowOrigins: allowOrigins,
		YAML:         yaml,
	})
}

/* private */

var action = new(sync.Once)
var config *Config

// getEnv retrieves the value of the environment variable named by envkey or returns placeholder.
func getEnv(envkey, placeholder string) (value string) {
	value = os.Getenv(envkey)
	if len(value) == 0 {
		value = placeholder
	}
	return
}

func unmarshalYAML(bytes []byte) (result map[string]interface{}) {
	err := yaml.Unmarshal(bytes, &result)
	if err != nil {
		panic(err)
	}
	return
}
