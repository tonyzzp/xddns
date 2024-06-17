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
	app := &cli.App{
		Flags: []cli.Flag{flagConfig},
		Commands: []*cli.Command{
			cmdSet,
			cmdUpdate,
			{
				Name:   "test",
				Action: test,
			},
		},
		Before: func(ctx *cli.Context) error {
			configFile := ctx.String("config")
			log.Println("configFile", configFile)
			return config.Init(configFile)
		},
	}
	e := app.Run(os.Args)
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}
