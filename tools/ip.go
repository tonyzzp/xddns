package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func get(t string) (string, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, e := client.Get(fmt.Sprintf("https://%s.jsonip.com", t))
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

// 获取本机的公网ipv6地址。优先获取eui64格式地址。
// 如果本机没有公网地址，则会返回实际访问外网使用的ipv6地址(一般是上级路由公网地址),同时返回一个error
func GetExternalIpv6() (string, error) {
	remoteIP, e := get("ipv6")
	log.Println("jsonip", remoteIP)
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
	log.Println("所有公网全球单播ipv6数量", len(ips))
	for _, v := range ips {
		log.Println(v)
	}
	var found string
	if len(ips) > 0 {
		for _, ip := range ips {
			if isEUI64(ip.IP) {
				log.Println("找到eui64格式ip")
				found = ip.IP.String()
				break
			}
		}
		if found == "" {
			found = ips[0].IP.String()
		}
	}
	log.Println("found", found)
	if found != "" {
		return found, nil
	} else {
		return remoteIP, errors.New("no global unicast ipv6")
	}
}
