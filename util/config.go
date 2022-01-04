package util

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server Server
	Db     Db
}

type Server struct {
	Port string `yaml:"port"`
}

type Db struct {
	Path string `yaml:"path"`
}

var Configs = Config{}

func InitSetting() {
	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal("fail to read file:", err)
	}

	err = yaml.Unmarshal(file, &Configs)
	if err != nil {
		log.Fatal("fail to yaml unmarshal:", err)
	}
}
