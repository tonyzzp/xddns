package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func cmdDumpAction(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	resp, e := obtainClient(domain).ListAllRecords(domain)
	if e != nil {
		return e
	}
	fmt.Println("totalcount", len(resp))
	fmt.Printf("%50s %7s %50s\n", "domain", "type", "value")
	fmt.Println(strings.Repeat("-", 50+7+50+2))
	for _, record := range resp {
		value := record.Value
		if len(value) > 40 {
			value = value[:40]
		}
		fmt.Printf("%50s %7s %50s\n", record.Domain, record.Type, value)
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
