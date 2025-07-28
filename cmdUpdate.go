package main

import (
	"errors"
	"fmt"

	"github.com/tonyzzp/xddns/dns"
	"github.com/tonyzzp/xddns/tools"

	"github.com/urfave/cli/v2"
)

var cmdUpdate = &cli.Command{
	Name:   "update",
	Usage:  "设置为本机ip",
	Action: cmdUpdateAction,
	Flags: []cli.Flag{
		flagIpType,
		flagDomain,
		flagRegister,
	},
}

func cmdUpdateAction(ctx *cli.Context) error {
	register := flagRegister.Get(ctx)
	var domain = flagDomain.Get(ctx)
	var ipType = flagIpType.Get(ctx)
	var ip = getLocalIP(ipType)
	if ip == "" {
		return errors.New("无法获取ip")
	}
	fmt.Println("updateAction", domain, ipType, ip)
	var recordType = ""
	if ipType == "ipv4" {
		recordType = dns.RECORD_TYPE_A
	} else {
		recordType = dns.RECORD_TYPE_AAAA
	}
	client, e := obtainClient(register, domain)
	if e != nil {
		return e
	}
	return client.EditRecord(dns.EditRecordParams{
		Domain: domain,
		Type:   recordType,
		Value:  ip,
	})
}

func getLocalIP(ipType string) string {
	retryTimes := 0
	for retryTimes < 3 {
		retryTimes++
		var ip string
		var e error
		if ipType == "ipv4" {
			ip, e = tools.GetExternalIpv4()
		} else {
			ip, e = tools.GetExternalIpv6()
		}
		if e == nil && ip != "" {
			return ip
		}
	}
	return ""
}
