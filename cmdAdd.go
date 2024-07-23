package main

import (
	"fmt"

	"github.com/tonyzzp/xddns/dns"

	"github.com/urfave/cli/v2"
)

var cmdAdd = &cli.Command{
	Name:  "add",
	Usage: "添加一条新的记录",
	Flags: []cli.Flag{
		flagDomain,
		flagRecordType,
		flagValue,
	},
	Action: actionAdd,
}

func actionAdd(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	t := flagRecordType.Get(ctx)
	value := flagValue.Get(ctx)
	fmt.Println("actionAdd", domain, t, value)
	return obtainClient(domain).AddRecord(dns.AddRecordParams{
		Domain: domain,
		Type:   t,
		Value:  value,
	})
}
