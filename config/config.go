package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type appConfig struct {
	Region    string `yaml:"region"`
	KeyId     string `yaml:"keyid"`
	KeySecret string `yaml:"keysecret"`
}

var Config *appConfig

func Init(configFile string) error {
	Config = &appConfig{}
	r, e := os.OpenFile(configFile, os.O_RDONLY, os.ModePerm)
	if e != nil {
		return e
	}
	e = yaml.NewDecoder(r).Decode(Config)
	r.Close()
	return e
}
