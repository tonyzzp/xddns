package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"xddns/dns"

	"github.com/urfave/cli/v2"
)

var cmdUpdate = &cli.Command{
	Name:   "update",
	Usage:  "设置为本机ip",
	Action: cmdUpdateAction,
	Flags: []cli.Flag{
		flagIpType,
		flagDomain,
	},
}

func cmdUpdateAction(ctx *cli.Context) error {
	var domain = flagDomain.Get(ctx)
	var ipType = flagIpType.Get(ctx)
	var ip = getLocalIP(ipType)
	if ip == "" {
		return errors.New("无法获取ip")
	}
	log.Println("updateAction", domain, ipType, ip)
	var recordType = ""
	if ipType == "ipv4" {
		recordType = dns.RECORD_TYPE_A
	} else {
		recordType = dns.RECORD_TYPE_AAAA
	}
	return obtainClient(domain).EditRecord(dns.EditRecordParams{
		Domain: domain,
		Type:   recordType,
		Value:  ip,
	})
}

func getLocalIP(ipType string) string {
	url := fmt.Sprintf("https://%s.jsonip.com", ipType)
	var fetch = func() (string, error) {
		resp, e := http.Get(url)
		if e != nil {
			return "", e
		}
		if resp.StatusCode != 200 {
			return "", errors.New("net error")
		}
		m := map[string]any{}
		e = json.NewDecoder(resp.Body).Decode(&m)
		if e != nil {
			return "", e
		}
		aip, ok := m["ip"]
		if !ok {
			return "", errors.New("no ip")
		}
		return aip.(string), nil
	}
	retryTimes := 0
	for retryTimes < 3 {
		retryTimes++
		ip, e := fetch()
		if e == nil && ip != "" {
			return ip
		}
	}
	return ""
}
