package main

import (
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
	log.Println(c.ListMainDomains())
	log.Println(c.QueryRecords(dns.QueryRecordParams{Domain: "izzp.me"}))
	log.Println(c.QueryRecords(dns.QueryRecordParams{Domain: "www.izzp.me"}))
	log.Println(c.ListAllRecords("izzp.me"))

	log.Println("增加")
	log.Println(c.AddRecord(dns.AddRecordParams{
		Domain: "test.izzp.me",
		Type:   "A",
		Value:  "1.1.1.4",
	}))

	// log.Println("edit")
	// log.Println(c.EditRecord(dns.EditRecordParams{
	// 	Domain: "test.izzp.me",
	// 	Type:   "A",
	// 	Value:  "1.2.3.4",
	// }))

	log.Println("delete")
	log.Println(c.DelRecord(dns.DelRecordParams{
		Domain: "test.izzp.me",
		Type:   "A",
	}))
	return nil
}

func testAli(ctx *cli.Context) error {

	c := ali.New()
	{
		log.Println("mainDomains")
		list, e := c.ListMainDomains()
		log.Println(e)
		for _, v := range list {
			log.Println(v)
		}
	}

	{
		log.Println("listAll")
		list, e := c.ListAllRecords("veikr.com")
		log.Println(e)
		for _, v := range list {
			log.Println(v)
		}
	}

	{
		log.Println("find")
		list, e := c.QueryRecords(dns.QueryRecordParams{Domain: "veikr.com"})
		log.Println(e)
		for _, v := range list {
			log.Println(v)
		}
	}

	{
		log.Println("edit")
		e := c.EditRecord(dns.EditRecordParams{Domain: "test.veikr.com", Type: "CNAME", Value: "www.veikr.com"})
		log.Println(e)
	}

	{
		log.Println("add")
		e := c.AddRecord(dns.AddRecordParams{Domain: "test.veikr.com", Type: "CNAME", Value: "veikr.com"})
		log.Println(e)
	}

	{
		log.Println("del")
		e := c.DelRecord(dns.DelRecordParams{Domain: "test.veikr.com"})
		log.Println(e)
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
	log.SetOutput(os.Stdout)
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
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}
