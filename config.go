package goconf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

//config file path and ext,Ext default .json
type Config struct {
	Path  string
	Ext   string
	cache map[string]gjson.Result
	file  string
}

//if key does not exist, return error
//key: dev.user.name  dev is filename in config path
func (c *Config) Get(key string) (*gjson.Result, error) {
	keys := strings.Split(key, ".")

	if len(keys) < 2 {
		return nil, errors.New("config XPath is at least two paragraphs")
	}

	c.file = c.Path + keys[0] + c.Ext

	if _, err := os.Stat(c.file); os.IsNotExist(err) {
		return nil, errors.New("config path not found:" + c.file)
	}

	var result gjson.Result
	if c.cache[c.file].IsObject() {
		result = c.cache[c.file]
	} else {
		b, err := ioutil.ReadFile(c.file)
		if err != nil {
			return nil, errors.Wrap(err, "config file read err")
		}
		result = gjson.ParseBytes(b)
		c.cache[c.file] = result
	}

	ret := result.Get(strings.Join(keys[1:], "."))
	return &ret, nil
}

//Ignore the error and return a zero value when it does not exist
func (c *Config) MustGet(key string) *gjson.Result {
	ret, err := c.Get(key)
	if err != nil {
		return &gjson.Result{}
	}

	return ret
}

func NewConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("config path not found:" + path)
	}
	configPath, _ := filepath.Abs(path)
	config := &Config{
		Path:  configPath + "/",
		Ext:   ".json",
		cache: make(map[string]gjson.Result, 0),
	}
	return config, nil
}
