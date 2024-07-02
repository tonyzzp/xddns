package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func get(t string) (string, error) {
	resp, e := http.Get(fmt.Sprintf("https://%s.jsonip.com", t))
	if e != nil {
		return "", e
	}
	if resp.StatusCode != 200 {
		return "", errors.New("net error")
	}
	bs, e := io.ReadAll(resp.Body)
	if e != nil {
		return "", e
	}
	var m = make(map[string]any)
	e = json.Unmarshal(bs, &m)
	if e != nil {
		return "", e
	}
	var ip, ok = m["ip"]
	if !ok {
		return "", errors.New("fetch ip failed")
	}
	return ip.(string), nil
}

func isEUI64(ip net.IP) bool {
	if ip.To4() != nil {
		return false
	}
	return ip[11] == 0xff && ip[12] == 0xfe
}

func GetExternalIpv4() (string, error) {
	return get("ipv4")
}

func GetExternalIpv6() (string, error) {
	ip, e := get("ipv6")
	log.Println("jsonip", ip)
	if e != nil {
		return "", e
	}
	interfaces, e := net.Interfaces()
	if e != nil {
		return "", e
	}
	var ips = make([]*net.IPNet, 0)
	for _, inter := range interfaces {
		addrs, e := inter.Addrs()
		if e != nil {
			continue
		}
		log.Println("interface", inter)
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				{
					ones, _ := v.Mask.Size()
					log.Println(
						"network", v.Network(),
						"ip", v.IP,
						"IsGlobalUnicast", v.IP.IsGlobalUnicast(),
						"IsPrivate", v.IP.IsPrivate(),
						"IsLinkLocalUnicast", v.IP.IsLinkLocalUnicast(),
						"mask", ones,
					)
					if v.IP.To4() == nil && v.IP.IsGlobalUnicast() && !v.IP.IsPrivate() {
						ips = append(ips, v)
					}
					break
				}
			}
		}
	}
	log.Println("所有全球单播ipv6")
	for _, v := range ips {
		log.Println(v)
		if isEUI64(v.IP) {
			return v.IP.String(), nil
		}
	}
	if len(ips) > 0 {
		return ips[0].String(), nil
	}
	return ip, nil
}
