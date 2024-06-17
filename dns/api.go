package dns

import (
	"ali-ddns/config"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

const RECORD_TYPE_A = "A"
const RECORD_TYPE_AAAA = "AAAA"
const RECORD_TYPE_CNAME = "CNAME"
const RECORD_TYPE_TXT = "TXT"

var dnsClient *alidns.Client

func initDnsClient() error {
	var e error
	if dnsClient == nil {
		dnsClient, e = alidns.NewClientWithAccessKey(config.Config.Region, config.Config.KeyId, config.Config.KeySecret)
	}
	return e
}

func QueryRecords(domain string) ([]alidns.Record, error) {
	e := initDnsClient()
	if e != nil {
		return nil, e
	}
	r := alidns.CreateDescribeSubDomainRecordsRequest()
	r.SubDomain = domain
	res, e := dnsClient.DescribeSubDomainRecords(r)
	if e != nil {
		return nil, e
	}
	return res.DomainRecords.Record, nil
}

func QueryRecord(subDomain string, recordType string) (*alidns.Record, error) {
	e := initDnsClient()
	if e != nil {
		return nil, e
	}
	log.Println("dns.QueryRecord", subDomain, recordType)
	r := alidns.CreateDescribeSubDomainRecordsRequest()
	r.SubDomain = subDomain
	r.Type = recordType
	res, e := dnsClient.DescribeSubDomainRecords(r)
	if e != nil {
		log.Println("dns.QueryRecord failed", e.Error())
		return nil, e
	}
	if len(res.DomainRecords.Record) > 0 {
		record := &res.DomainRecords.Record[0]
		log.Println("dns.QueryRecord result", record)
		return record, nil
	}
	return nil, nil
}

func AddRecord(domain string, rr string, recordType string, value string) error {
	e := initDnsClient()
	if e != nil {
		return nil
	}
	log.Println("dns.AddRecord", domain, rr, recordType, value)
	req := alidns.CreateAddDomainRecordRequest()
	req.DomainName = domain
	req.RR = rr
	req.Type = recordType
	req.Value = value
	_, e = dnsClient.AddDomainRecord(req)
	log.Println("dns.AddRecord result", e)
	return e
}

func UpdateRecord(recordId string, rr string, recordType string, value string) error {
	e := initDnsClient()
	if e != nil {
		return nil
	}
	log.Println("dns.UpdateRecord", recordId, rr, recordType, value)
	req := alidns.CreateUpdateDomainRecordRequest()
	req.RecordId = recordId
	req.RR = rr
	req.Type = recordType
	req.Value = value
	_, e = dnsClient.UpdateDomainRecord(req)
	se := e.(*errors.ServerError)
	if se != nil && se.ErrorCode() == "DomainRecordDuplicate" {
		log.Println("dns.UpdateRecord 成功。相同记录已存在")
		return nil
	}
	log.Println("dns.UpdateRecord result", e)
	return e
}

func EditRecord(domain string, rr string, recordType string, value string) error {
	exist, e := QueryRecord(rr+"."+domain, recordType)
	if e != nil {
		return e
	}
	if exist == nil {
		return AddRecord(domain, rr, recordType, value)
	} else {
		return UpdateRecord(exist.RecordId, rr, recordType, value)
	}
}
