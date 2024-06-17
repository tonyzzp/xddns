package main

import "github.com/urfave/cli/v2"

var flagRR = &cli.StringFlag{
	Name:     "value",
	Required: true,
}

var flagDomain = &cli.StringFlag{
	Name:     "domain",
	Required: true,
}

var flagRecordType = &cli.StringFlag{
	Name:     "type",
	Required: true,
}

var flagValue = &cli.StringFlag{
	Name:     "value",
	Required: true,
}

var flagDomains = &cli.StringFlag{
	Name:     "domains",
	Usage:    "域名。多个使用逗号隔开",
	Required: true,
}

var flagIpType = &cli.StringFlag{
	Name:     "type",
	Usage:    "ip类型:  ipv4/ipv6",
	Required: true,
}
