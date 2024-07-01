package main

import (
	"ali-ddns/dns"

	"github.com/urfave/cli/v2"
)

func cmdSetAction(ctx *cli.Context) error {
	t := flagRecordType.Get(ctx)
	domain := flagDomain.Get(ctx)
	value := flagValue.Get(ctx)
	return dns.EditRecord(domain, t, value)
}

var cmdSet = &cli.Command{
	Name:  "set",
	Usage: "set dns record",
	Flags: []cli.Flag{
		flagDomain,
		flagRecordType,
		flagValue,
	},
	Action: cmdSetAction,
}
