package main

import (
	"ali-ddns/config"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func test(ctx *cli.Context) error {
	return nil
}

func findConfigFile() string {
	name := "ali_config.yaml"
	fi, e := os.Stat(name)
	if e == nil && fi != nil && !fi.IsDir() {
		return name
	}
	dir, e := os.Executable()
	if e != nil {
		return name
	}
	file := filepath.Join(filepath.Dir(dir), name)
	return file
}

func main() {
	e := config.Init(findConfigFile())
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			cmdSet,
			cmdUpdate,
			{
				Name:   "test",
				Action: test,
			},
		},
	}
	e = app.Run(os.Args)
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}
