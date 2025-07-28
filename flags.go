package main

import "github.com/urfave/cli/v2"

var flagDomain = &cli.StringFlag{
	Name:     "domain",
	Required: true,
	Usage:    "完整域名，如: mydomain.com, music.mydomain.com",
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

var flagTTL = &cli.IntFlag{
	Name:     "ttl",
	Usage:    "域名缓存时长(秒)",
	Required: false,
}

var flagRegister = &cli.StringFlag{
	Name:     "register",
	Usage:    "域名注册商 ali/cf。如果不传会自动判断，如果某个域名在多个服务商处都配置有解析，可手动传此参数",
	Required: false,
}
