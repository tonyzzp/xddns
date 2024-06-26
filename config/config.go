package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type appConfig struct {
	Region    string `yaml:"region"`
	KeyId     string `yaml:"keyid"`
	KeySecret string `yaml:"keysecret"`
}

var Config *appConfig

func findConfigFile(file string) string {
	if file != "" {
		return file
	}

	// 查找工作目录
	name := "ali_config.yaml"
	fi, e := os.Stat(name)
	if e == nil && fi != nil && !fi.IsDir() {
		return name
	}

	// 查找exe目录
	dir, e := os.Executable()
	if e != nil {
		return name
	}
	file = filepath.Join(filepath.Dir(dir), name)
	return file
}

func Init(file string) error {
	file = findConfigFile(file)
	Config = &appConfig{}
	r, e := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if e != nil {
		return e
	}
	e = yaml.NewDecoder(r).Decode(Config)
	r.Close()
	return e
}
