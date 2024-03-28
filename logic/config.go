package logic

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ApiKey   string `yaml:"apiKey"`
	Port     string `yaml:"port"`
}

var Data Config

func LoadConfig() error {
	// 读取配置文件
	configData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	// 解析配置文件
	err = yaml.Unmarshal(configData, &Data)
	if err != nil {
		return err
	}

	return nil
}
