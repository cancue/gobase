package util

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"gopkg.in/yaml.v3"
)

// GetEnv retrieves the value of the environment variable named by the key with placeholder.
func GetEnv(key string, placeholder string) (value string) {
	value = os.Getenv(key)
	if len(value) == 0 {
		value = placeholder
	}

	return
}

// CallerDir retrives the path of caller for this method.
func CallerDir() string {
	_, b, _, _ := runtime.Caller(0)

	return path.Dir(b)
}

// ReadYAMLFile returns a yaml file corresponding to the path in an operable form.
func ReadYAMLFile(path string) (result map[string]interface{}) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return readYAML(file)
}

func readYAML(file io.Reader) (result map[string]interface{}) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		panic(err)
	}

	return
}
