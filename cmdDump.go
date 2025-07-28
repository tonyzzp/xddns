package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func cmdDumpAction(ctx *cli.Context) error {
	register := flagRegister.Get(ctx)
	domain := flagDomain.Get(ctx)
	client, e := obtainClient(register, domain)
	if e != nil {
		return e
	}
	resp, e := client.ListAllRecords(domain)
	if e != nil {
		return e
	}
	fmt.Println("totalcount", len(resp))
	fmt.Printf("%50s %7s %7s %50s %6s\n", "domain", "type", "proxied", "value", "tlt")
	fmt.Println(strings.Repeat("-", 130))
	for _, record := range resp {
		value := record.Value
		if len(value) > 40 {
			value = value[:40]
		}
		var proxied = ""
		if record.Proxied {
			proxied = "proxied"
		}
		fmt.Printf("%50s %7s %7s %50s %6d\n", record.Domain, record.Type, proxied, value, record.TLT)
	}
	return nil
}

var cmdDump = &cli.Command{
	Name:  "dump",
	Usage: "dump all dns records",
	Flags: []cli.Flag{
		flagDomain,
		flagRegister,
	},
	Action: cmdDumpAction,
}
