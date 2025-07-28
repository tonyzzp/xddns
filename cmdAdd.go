package main

import (
	"errors"
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
		flagTTL,
		flagRegister,
	},
	Action: actionAdd,
}

func actionAdd(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	t := flagRecordType.Get(ctx)
	value := flagValue.Get(ctx)
	ttl := flagTTL.Get(ctx)
	register := flagRegister.Get(ctx)
	if t == "" {
		return errors.New("need param --type")
	}
	fmt.Println("actionAdd", domain, t, value)
	client, e := obtainClient(register, domain)
	if e != nil {
		return e
	}
	return client.AddRecord(dns.AddRecordParams{
		Domain: domain,
		Type:   t,
		Value:  value,
		TTL:    ttl,
	})
}
