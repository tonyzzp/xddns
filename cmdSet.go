package main

import (
	"fmt"

	"github.com/tonyzzp/xddns/dns"

	"github.com/urfave/cli/v2"
)

func cmdSetAction(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	t := flagRecordType.Get(ctx)
	value := flagValue.Get(ctx)
	fmt.Println("actionSet", domain, t, value)
	return obtainClient(domain).EditRecord(dns.EditRecordParams{
		Domain: domain,
		Type:   t,
		Value:  value,
	})
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
