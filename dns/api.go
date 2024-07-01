package dns

import (
	"ali-ddns/config"
	"errors"
	"log"
	"strings"

	alierrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

const RECORD_TYPE_A = "A"
const RECORD_TYPE_AAAA = "AAAA"
const RECORD_TYPE_CNAME = "CNAME"
const RECORD_TYPE_TXT = "TXT"

type DomainResolved struct {
	DomainName string
	RR         string
}

type AddRecordParams struct {
	Domain string
	Type   string
	Value  string
}

type UpdateRecordParams struct {
	Id     string
	Domain string
	Type   string
	Value  string
}

type DelRecordParams struct {
	Domain string
	Type   string
}

type EditRecordParams struct {
	Domain string
	Type   string
	Value  string
}

var dnsClient *alidns.Client

func initDnsClient() error {
	var e error
	if dnsClient == nil {
		dnsClient, e = alidns.NewClientWithAccessKey(config.Config.Region, config.Config.KeyId, config.Config.KeySecret)
	}
	return e
}

func QueryDomains() ([]string, error) {
	log.Println("dns.QueryDomains")
	e := initDnsClient()
	if e != nil {
		return nil, e
	}
	var rtn []string
	req := alidns.CreateDescribeDomainsRequest()
	res, e := dnsClient.DescribeDomains(req)
	if e == nil && res.IsSuccess() {
		rtn = make([]string, 0)
		for _, domain := range res.Domains.Domain {
			rtn = append(rtn, domain.DomainName)
		}
	}
	log.Println("dns.QueryDomains.result", rtn)
	return rtn, e
}

func ResolveDomain(fullDomain string) (*DomainResolved, error) {
	log.Println("dns.ResolveDomain", fullDomain)
	domains, e := QueryDomains()
	if e != nil {
		return nil, e
	}
	for _, domain := range domains {
		if strings.HasSuffix(fullDomain, domain) {
			rr := strings.TrimSuffix(fullDomain, domain)
			if strings.HasSuffix(rr, ".") {
				rr = rr[:len(rr)-1]
			} else if rr == "" {
				rr = "@"
			}
			rtn := &DomainResolved{
				DomainName: domain,
				RR:         rr,
			}
			log.Println("dns.ResolveDomain.result", *rtn)
			return rtn, nil
		}
	}
	return nil, errors.New("not exists")
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

func QueryRecord(domain string, recordType string) (*alidns.Record, error) {
	e := initDnsClient()
	if e != nil {
		return nil, e
	}
	log.Println("dns.QueryRecord", domain, recordType)
	r := alidns.CreateDescribeSubDomainRecordsRequest()
	r.SubDomain = domain
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

func AddRecord(params AddRecordParams) error {
	e := initDnsClient()
	if e != nil {
		return nil
	}
	log.Println("dns.AddRecord", params)
	info, e := ResolveDomain(params.Domain)
	if e != nil {
		return e
	}
	req := alidns.CreateAddDomainRecordRequest()
	req.DomainName = info.DomainName
	req.RR = info.RR
	req.Type = params.Type
	req.Value = params.Value
	_, e = dnsClient.AddDomainRecord(req)
	log.Println("dns.AddRecord result", e)
	return e
}

func UpdateRecord(params UpdateRecordParams) error {
	e := initDnsClient()
	if e != nil {
		return nil
	}
	info, e := ResolveDomain(params.Domain)
	if e != nil {
		return e
	}
	log.Println("dns.UpdateRecord", params)
	req := alidns.CreateUpdateDomainRecordRequest()
	req.RecordId = params.Id
	req.Type = params.Type
	req.Value = params.Value
	req.RR = info.RR
	_, e = dnsClient.UpdateDomainRecord(req)
	se, _ := e.(*alierrors.ServerError)
	if se != nil && se.ErrorCode() == "DomainRecordDuplicate" {
		log.Println("dns.UpdateRecord 成功。相同记录已存在")
		return nil
	}
	log.Println("dns.UpdateRecord result", e)
	return e
}

func EditRecord(params EditRecordParams) error {
	exist, e := QueryRecord(params.Domain, params.Type)
	if e != nil {
		return e
	}
	if exist == nil {
		return AddRecord(AddRecordParams{
			Domain: params.Domain,
			Type:   params.Type,
			Value:  params.Value,
		})
	} else {
		return UpdateRecord(UpdateRecordParams{
			Id:     exist.RecordId,
			Domain: params.Domain,
			Type:   params.Type,
			Value:  params.Value,
		})
	}
}

func DelRecord(params DelRecordParams) error {
	log.Println("dns.DelRecord", params)
	e := initDnsClient()
	if e != nil {
		return nil
	}
	exist, e := QueryRecord(params.Domain, params.Type)
	if e != nil {
		return e
	}
	if exist == nil {
		log.Println("记录不存在")
		return nil
	}
	req := alidns.CreateDeleteSubDomainRecordsRequest()
	req.DomainName = exist.DomainName
	req.RR = exist.RR
	req.Type = exist.Type
	resp, e := dnsClient.DeleteSubDomainRecords(req)
	if e != nil {
		return e
	}
	if resp.IsSuccess() {
		return nil
	}
	return errors.New("delrecord failed: " + resp.String())
}

func GetAllRecords(domain string) (*alidns.DescribeDomainRecordsResponse, error) {
	log.Println("dns.GetAllRecords", domain)
	e := initDnsClient()
	if e != nil {
		return nil, e
	}

	req := alidns.CreateDescribeDomainRecordsRequest()
	req.DomainName = domain
	req.PageSize = "100"
	resp, e := dnsClient.DescribeDomainRecords(req)
	return resp, e
}
