package main

import (
	"ali-ddns/ali"
	"log"

	"github.com/urfave/cli/v2"
)

func cmdDumpAction(ctx *cli.Context) error {
	domain := flagDomain.Get(ctx)
	resp, e := ali.GetAllRecords(domain)
	if e != nil {
		return e
	}
	log.Println("totalcount", resp.TotalCount)
	log.Println("records:", len(resp.DomainRecords.Record))
	for _, record := range resp.DomainRecords.Record {
		value := record.Value
		if len(value) > 40 {
			value = value[:40]
		}
		log.Printf("%7s %40s %7s %40s %04d %02d", record.Status, record.RR, record.Type, value, record.TTL, record.Weight)
	}
	return nil
}

var cmdDump = &cli.Command{
	Name:  "dump",
	Usage: "dump all dns records",
	Flags: []cli.Flag{
		flagDomain,
	},
	Action: cmdDumpAction,
}
