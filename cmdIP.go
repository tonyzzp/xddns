package main

import (
	"fmt"

	"github.com/tonyzzp/xddns/tools"

	"github.com/urfave/cli/v2"
)

func actionIP(ctx *cli.Context) error {
	ipv4, e := tools.GetExternalIpv4()
	fmt.Println("ipv4", ipv4, e)

	ipv6, e := tools.GetExternalIpv6()
	fmt.Println("ipv6", ipv6, e)
	return nil
}

var cmdIP = &cli.Command{
	Name:   "ip",
	Usage:  "print my external ips",
	Action: actionIP,
}
