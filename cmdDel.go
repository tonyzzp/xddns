package main

import (
	"ali-ddns/dns"

	"github.com/urfave/cli/v2"
)

func cmdDelAction(ctx *cli.Context) error {
	t := ctx.String("type")
	domain := ctx.String("domain")
	rr := ctx.String("rr")
	return dns.DelRecord(domain, rr, t)
}

var cmdDel = &cli.Command{
	Name:  "del",
	Usage: "del dns record",
	Flags: []cli.Flag{
		flagDomain,
		flagRR,
		flagRecordType,
	},
	Action: cmdDelAction,
}
