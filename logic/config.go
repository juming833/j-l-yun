package logic

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	ApiKey    string `yaml:"apiKey"`
	Port      string `yaml:"port"`
	Token     string `yaml:"token"`
	Loglevel  string `yaml:"loglevel"`
	Test      bool   `yaml:"test"`
	CacheTime int    `yaml:"cache_time"`
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
