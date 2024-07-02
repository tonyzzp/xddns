package main

import (
	"log"

	"github.com/urfave/cli/v2"
)

func cmdDumpAction(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	resp, e := obtainClient(domain).ListAllRecords(domain)
	if e != nil {
		return e
	}
	log.Println("totalcount", len(resp))
	for _, record := range resp {
		value := record.Value
		if len(value) > 40 {
			value = value[:40]
		}
		log.Printf("%50s %7s %50s", record.Domain, record.Type, value)
	}
	return nil
}

var cmdDump = &cli.Command{
	Name:  "dump",
	Usage: "dump all dns records",
	Flags: []cli.Flag{
		flagDomain,
	},
	Action: cmdDumpAction,
}
