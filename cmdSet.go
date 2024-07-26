package main

import (
	"errors"
	"fmt"

	"github.com/tonyzzp/xddns/dns"

	"github.com/urfave/cli/v2"
)

func cmdSetAction(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	t := flagRecordType.Get(ctx)
	value := flagValue.Get(ctx)
	if t == "" {
		return errors.New("need param --type")
	}
	fmt.Println("actionSet", domain, t, value)
	client, e := obtainClient(domain)
	if e != nil {
		return e
	}
	return client.EditRecord(dns.EditRecordParams{
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
