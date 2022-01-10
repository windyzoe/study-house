package util

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server Server
	Db     Db
	Auth   Auth
}

type Server struct {
	Port string `yaml:"port"`
}

type Db struct {
	Path string `yaml:"path"`
}

type Auth struct {
	Whitelist []string `yaml:"whitelist"`
}

var Configs = Config{}

func InitConfig() {
	file, err := ioutil.ReadFile("./config-" + *FLAG_ENV + ".yaml")
	if err != nil {
		log.Error().Err(err).Msg("fail to read file:")
	}

	err = yaml.Unmarshal(file, &Configs)
	if err != nil {
		log.Error().Err(err).Msg("fail to yaml unmarshal:")

	}
	log.Printf("ENV::%s, CONFIG::%s", *FLAG_ENV, Configs)
}
