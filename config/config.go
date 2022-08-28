package config

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Cfg *Conf

type Conf struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	Api struct {
		Port int64 `yaml:"port"`
	}
	Contracts struct {
		Exchange  string `yaml:"exchange"`
		PooHub    string `yaml:"pool_hub"`
		Positions string `yaml:"positions"`
	}
	Chain struct {
		ChainID    int64  `yaml:"chain_id"`
		URL        string `yaml:"url"`
		PrivateKey string `yaml:"private_key"`
	}
}

func LoadConfig(configFile string) *Conf {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Errorf("fail to read config %v", err))
	}

	var conf Conf
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		logrus.Fatalf("fail to parse config yaml  %q: %v", configFile, err)
	}

	Cfg = &conf
	return &conf
}
