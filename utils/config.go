package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Config 配置文件
type Config struct {
	Source   string
	Encrypt  string
	Decrypt  string
	Download string
}

//ReadConfig 读取配置文件
func ReadConfig(source string) (config *Config, err error) {
	configBytes, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configBytes, &config)
	return config, err
}
