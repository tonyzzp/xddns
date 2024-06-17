package main

import (
	"ali-ddns/dns"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

var cmdUpdate = &cli.Command{
	Name:   "update",
	Usage:  "设置为本机ip",
	Action: cmdUpdateAction,
	Flags: []cli.Flag{
		flagIpType,
		flagDomains,
	},
}

func _update(domains []string, recordType string, value string) error {
	for _, v := range domains {
		strs := strings.Split(v, ".")
		var domain = ""
		var rr = ""
		if len(strs) == 2 {
			domain = v
			rr = "@"
		} else {
			domain = strings.Join(strs[len(strs)-2:], ".")
			rr = strings.Join(strs[0:len(strs)-2], ".")
		}
		log.Println("_update", domain, rr)
		e := dns.EditRecord(domain, rr, recordType, value)
		if e != nil {
			return e
		}
	}
	return nil
}

func cmdUpdateAction(ctx *cli.Context) error {
	var domains = strings.Split(ctx.String("domains"), ",")
	var ipType = ctx.String("type")
	var ip = getLocalIP(ipType)
	if ip == "" {
		return errors.New("无法获取ip")
	}
	log.Println("updateAction", domains, ipType, ip)
	var recordType = ""
	if ipType == "ipv4" {
		recordType = dns.RECORD_TYPE_A
	} else {
		recordType = dns.RECORD_TYPE_AAAA
	}
	return _update(domains, recordType, ip)
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
