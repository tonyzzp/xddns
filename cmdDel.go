package main

import (
	"fmt"

	"github.com/tonyzzp/xddns/dns"

	"github.com/urfave/cli/v2"
)

func cmdDelAction(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	t := flagRecordType.Get(ctx)
	fmt.Println("delAction", domain, t)
	return obtainClient(domain).DelRecord(dns.DelRecordParams{
		Domain: domain,
		Type:   t,
	})
}

var cmdDel = &cli.Command{
	Name:  "del",
	Usage: "del dns records",
	Flags: []cli.Flag{
		flagDomain,
		flagRecordType,
	},
	Action: cmdDelAction,
}
