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
	Key      string
	M3u8     string
	Keyurl   string
	Tsurl    string
	Finalurl string
}

type ResourceConfig struct {
	Cdn string
}

//ReadConfig 读取配置文件
func ReadConfig() (config *Config, err error) {
	configBytes, err := ioutil.ReadFile(".yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configBytes, &config)
	return config, err
}``

//ReadResourceConfig 读取配置文件
func ReadResourceConfig() (config *ResourceConfig, err error) {
	configBytes, err := ioutil.ReadFile(".yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configBytes, &config)
	return config, err
}
