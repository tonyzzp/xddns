package main

import (
	"fmt"

	"github.com/tonyzzp/xddns/dns"

	"github.com/urfave/cli/v2"
)

func cmdDelAction(ctx *cli.Context) error {
	register := flagRegister.Get(ctx)
	domain := flagDomain.Get(ctx)
	t := flagRecordType.Get(ctx)
	fmt.Println("delAction", domain, t)
	client, e := obtainClient(register, domain)
	if e != nil {
		return e
	}
	return client.DelRecord(dns.DelRecordParams{
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
		flagRegister,
	},
	Action: cmdDelAction,
}
