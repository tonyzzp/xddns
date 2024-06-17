package main

import (
	"ali-ddns/config"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func test(ctx *cli.Context) error {
	return nil
}

func main() {
	e := config.Init("ali_config.yaml")
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
