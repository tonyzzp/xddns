package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tonyzzp/xddns/ali"
	"github.com/tonyzzp/xddns/cf"
	"github.com/tonyzzp/xddns/config"
	"github.com/tonyzzp/xddns/dns"

	"github.com/urfave/cli/v2"
)

func testCloudFlare(ctx *cli.Context) error {
	c := cf.New()
	fmt.Println(c.ListMainDomains())
	fmt.Println(c.QueryRecords(dns.QueryRecordParams{Domain: "izzp.me"}))
	fmt.Println(c.QueryRecords(dns.QueryRecordParams{Domain: "www.izzp.me"}))
	fmt.Println(c.ListAllRecords("izzp.me"))

	fmt.Println("增加")
	fmt.Println(c.AddRecord(dns.AddRecordParams{
		Domain: "test.izzp.me",
		Type:   "A",
		Value:  "1.1.1.4",
	}))

	// fmt.Println("edit")
	// fmt.Println(c.EditRecord(dns.EditRecordParams{
	// 	Domain: "test.izzp.me",
	// 	Type:   "A",
	// 	Value:  "1.2.3.4",
	// }))

	fmt.Println("delete")
	fmt.Println(c.DelRecord(dns.DelRecordParams{
		Domain: "test.izzp.me",
		Type:   "A",
	}))
	return nil
}

func testAli(ctx *cli.Context) error {

	c := ali.New()
	{
		fmt.Println("mainDomains")
		list, e := c.ListMainDomains()
		fmt.Println(e)
		for _, v := range list {
			fmt.Println(v)
		}
	}

	{
		fmt.Println("listAll")
		list, e := c.ListAllRecords("veikr.com")
		fmt.Println(e)
		for _, v := range list {
			fmt.Println(v)
		}
	}

	{
		fmt.Println("find")
		list, e := c.QueryRecords(dns.QueryRecordParams{Domain: "veikr.com"})
		fmt.Println(e)
		for _, v := range list {
			fmt.Println(v)
		}
	}

	{
		fmt.Println("edit")
		e := c.EditRecord(dns.EditRecordParams{Domain: "test.veikr.com", Type: "CNAME", Value: "www.veikr.com"})
		fmt.Println(e)
	}

	{
		fmt.Println("add")
		e := c.AddRecord(dns.AddRecordParams{Domain: "test.veikr.com", Type: "CNAME", Value: "veikr.com"})
		fmt.Println(e)
	}

	{
		fmt.Println("del")
		e := c.DelRecord(dns.DelRecordParams{Domain: "test.veikr.com"})
		fmt.Println(e)
	}

	return nil
}

func findDnsServer(fullDomain string) string {
	for _, v := range config.Config.Ali.Domains {
		if strings.HasSuffix(fullDomain, v) {
			return "ali"
		}
	}
	for v, _ := range config.Config.CloudFlare.Domains {
		if strings.HasSuffix(fullDomain, v) {
			return "cloudflare"
		}
	}
	return ""
}

var _aliClient dns.IDns
var _cloudflareClient dns.IDns

func obtainClient(fullDomain string) dns.IDns {
	var server = findDnsServer(fullDomain)
	if server == "ali" {
		if _aliClient == nil {
			_aliClient = ali.New()
		}
		return _aliClient
	} else if server == "cloudflare" {
		if _cloudflareClient == nil {
			_cloudflareClient = cf.New()
		}
		return _cloudflareClient
	}
	return nil
}

func main() {
	logFile, _ := os.OpenFile("xddns.log", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	app := &cli.App{
		Usage: "操作阿里dns解析记录",
		Flags: []cli.Flag{flagConfig},
		Commands: []*cli.Command{
			cmdAdd,
			cmdSet,
			cmdUpdate,
			cmdDel,
			cmdDump,
			cmdIP,
			{
				Name:   "testCloudFlare",
				Action: testCloudFlare,
			},
			{
				Name:   "testAli",
				Action: testAli,
			},
		},
		Before: func(ctx *cli.Context) error {
			var actionName = ctx.Args().First()
			if actionName != "ip" && actionName != "" {
				configFile := ctx.String("config")
				log.Println("configFile", configFile)
				return config.Init(configFile)
			}
			return nil
		},
	}
	e := app.Run(os.Args)
	if e == nil {
		fmt.Println("execute success")
	} else {
		fmt.Println(e)
		os.Exit(1)
	}
}
