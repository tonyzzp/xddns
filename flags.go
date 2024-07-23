package main

import "github.com/urfave/cli/v2"

var flagDomain = &cli.StringFlag{
	Name:     "domain",
	Required: true,
	Usage:    "完整域名，如: veikr.com , qm.pi.veikr.com",
}

var flagRecordType = &cli.StringFlag{
	Name:     "type",
	Required: false,
	Usage:    "A, AAAA, CNAME, TXT",
}

var flagValue = &cli.StringFlag{
	Name:     "value",
	Required: true,
}

var flagIpType = &cli.StringFlag{
	Name:     "type",
	Usage:    "ip类型:  ipv4/ipv6",
	Required: true,
}

var flagConfig = &cli.StringFlag{
	Name:     "config",
	Usage:    "xddns-config.yaml文件",
	Required: false,
}
