package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Config 配置文件
type Config struct {
	Originm3u8   string
	Origints     string
	Querykey     string
	Keyurl       string
	Cdn          string
	M3u8rule     string
	Encrypttsurl string
}

var config *Config

//ReadConfig 读取配置文件
func ReadConfig() *Config {
	return config
}

func InitConfig() error {
	if configBytes, err := ioutil.ReadFile(".yaml"); err != nil {
		return err
	} else {
		err := yaml.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}
	}
	return nil
}
