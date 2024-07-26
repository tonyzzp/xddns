package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func cmdDumpAction(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	client, e := obtainClient(domain)
	if e != nil {
		return e
	}
	resp, e := client.ListAllRecords(domain)
	if e != nil {
		return e
	}
	fmt.Println("totalcount", len(resp))
	fmt.Printf("%50s %7s %7s %50s\n", "domain", "type", "proxied", "value")
	fmt.Println(strings.Repeat("-", 50+7+50+3))
	for _, record := range resp {
		value := record.Value
		if len(value) > 40 {
			value = value[:40]
		}
		var proxied = ""
		if record.Proxied {
			proxied = "proxied"
		}
		fmt.Printf("%50s %7s %7s %50s\n", record.Domain, record.Type, proxied, value)
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
