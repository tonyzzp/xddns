package main

import (
	"ali-ddns/dns"

	"github.com/urfave/cli/v2"
)

func cmdSetAction(ctx *cli.Context) error {
	t := ctx.String("type")
	domain := ctx.String("domain")
	rr := ctx.String("rr")
	value := ctx.String("value")
	return dns.EditRecord(domain, rr, t, value)
}

var cmdSet = &cli.Command{
	Name:  "set",
	Usage: "set dns record",
	Flags: []cli.Flag{
		flagDomain,
		flagRR,
		flagRecordType,
		flagValue,
	},
	Action: cmdSetAction,
}
