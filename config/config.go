/*
Package config has weakness because of lenCaller.

Gobase.Start must be called in the same path with the config folder.
It is a fragile assumption, but worth it. Because Gobase is made for personal simplicity.
*/
package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/cancue/gobase/util"
	"gopkg.in/yaml.v3"
)

var configOnce = new(sync.Once)
var config = new(Config)

// Config public
type Config struct {
	Stage string
	YAML  map[string]interface{}
}

// Get retrieves config
func Get() *Config {
	return config
}

// Set sets and retrieves Config from yaml files on config directory.
//
// relativeDepth format follows '.', '..', '../..'...
// The name of the yaml file is determined by an environment variable whose key is 'STAGE'.
// If there is no environment variable, 'local' is selected by default. (e.g. 'local.yaml')
func Set(realativeDepth string) *Config {
	return SetWithCallerLength(0, realativeDepth)
}

// SetWithCallerLength is same as Set but requires an extra caller length.
func SetWithCallerLength(callerLength int, relativeDepth string) *Config {
	configOnce.Do(func() {
		config.Stage = util.GetEnv("STAGE", "local")
		_, b, _, _ := runtime.Caller(5 + callerLength)
		path := filepath.Join(path.Dir(b), relativeDepth, "config", config.Stage+".yaml")
		config.YAML = readYAMLFileWithEnvs(path)
	})

	return config
}

func readYAMLFileWithEnvs(p string) (result map[string]interface{}) {
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	str := os.ExpandEnv(string(bytes))

	err = yaml.Unmarshal([]byte(str), &result)
	if err != nil {
		panic(err)
	}

	return
}
