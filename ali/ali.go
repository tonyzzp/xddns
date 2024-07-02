package ali

import (
	"errors"
	"log"
	"strings"
	"xddns/config"
	"xddns/dns"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type DnsAli struct {
	client *alidns.Client
	cfg    config.ConfigAli
}

func New() *DnsAli {
	rtn := &DnsAli{
		cfg: config.Config.Ali,
	}
	dnsClient, e := alidns.NewClientWithAccessKey(config.Config.Ali.Region, config.Config.Ali.KeyId, config.Config.Ali.KeySecret)
	if e != nil {
		return nil
	}
	rtn.client = dnsClient
	return rtn
}

func (da *DnsAli) resolve(domain string) (dns.DomainResolved, error) {
	for _, d := range da.cfg.Domains {
		if strings.HasSuffix(domain, d) {
			rr := strings.TrimSuffix(domain, d)
			if strings.HasSuffix(rr, ".") {
				rr = rr[0 : len(rr)-1]
			} else if rr == "" {
				rr = "@"
			}
			return dns.DomainResolved{
				DomainName: d,
				RR:         rr,
			}, nil
		}
	}
	return dns.DomainResolved{}, errors.New("not exist")
}

func (da *DnsAli) ListMainDomains() ([]string, error) {
	return append([]string{}, da.cfg.Domains...), nil
}

func (da *DnsAli) ListAllRecords(mainDomain string) ([]dns.Record, error) {
	info, e := da.resolve(mainDomain)
	if e != nil {
		return nil, e
	}
	var page = 1
	doReq := func(page int) (int64, []dns.Record, error) {
		req := alidns.CreateDescribeDomainRecordsRequest()
		req.DomainName = info.DomainName
		req.PageNumber = requests.NewInteger(page)
		req.PageSize = requests.NewInteger(50)
		res, e := da.client.DescribeDomainRecords(req)
		if e != nil {
			return 0, nil, e
		}
		rtn := make([]dns.Record, 0)
		for _, v := range res.DomainRecords.Record {
			fullDomain := ""
			if v.RR == "" || v.RR == "@" {
				fullDomain = mainDomain
			} else {
				fullDomain = v.RR + "." + mainDomain
			}
			rtn = append(rtn, dns.Record{
				Id:      v.RecordId,
				Type:    v.Type,
				Enabled: v.Status == "ENABLE",
				Value:   v.Value,
				Domain:  fullDomain,
			})
		}
		return res.TotalCount, rtn, nil
	}
	rtn := make([]dns.Record, 0)
	for {
		total, list, e := doReq(page)
		if e != nil {
			return nil, e
		}
		rtn = append(rtn, list...)
		if total == int64(len(rtn)) {
			return rtn, nil
		}
		page++
	}
}

func (da *DnsAli) QueryRecords(params dns.QueryRecordParams) ([]dns.Record, error) {
	info, e := da.resolve(params.Domain)
	if e != nil {
		return nil, e
	}
	var page = 1
	doReq := func(page int) (int64, []dns.Record, error) {
		req := alidns.CreateDescribeSubDomainRecordsRequest()
		req.SubDomain = params.Domain
		req.Type = params.Type
		req.PageNumber = requests.NewInteger(page)
		req.PageSize = requests.NewInteger(50)
		res, e := da.client.DescribeSubDomainRecords(req)
		if e != nil {
			return 0, nil, e
		}
		rtn := make([]dns.Record, 0)
		for _, v := range res.DomainRecords.Record {
			fullDomain := ""
			if v.RR == "" || v.RR == "@" {
				fullDomain = info.DomainName
			} else {
				fullDomain = v.RR + "." + info.DomainName
			}
			rtn = append(rtn, dns.Record{
				Id:      v.RecordId,
				Type:    v.Type,
				Enabled: v.Status == "Enable",
				Value:   v.Value,
				Domain:  fullDomain,
			})
		}
		return res.TotalCount, rtn, nil
	}
	rtn := make([]dns.Record, 0)
	for {
		total, list, e := doReq(page)
		if e != nil {
			return nil, e
		}
		rtn = append(rtn, list...)
		if total == int64(len(rtn)) {
			return rtn, nil
		}
	}
}

func (da *DnsAli) AddRecord(params dns.AddRecordParams) error {
	info, e := da.resolve(params.Domain)
	if e != nil {
		return e
	}
	req := alidns.CreateAddDomainRecordRequest()
	req.DomainName = info.DomainName
	req.RR = info.RR
	req.Type = params.Type
	req.Value = params.Value
	_, e = da.client.AddDomainRecord(req)
	log.Println("dns.AddRecord result", e)
	return e
}

func (da *DnsAli) EditRecord(params dns.EditRecordParams) error {
	exists, e := da.QueryRecords(dns.QueryRecordParams{
		Domain: params.Domain,
		Type:   params.Type,
	})
	if e != nil {
		return e
	}

	if len(exists) > 0 {
		for _, record := range exists {
			req := alidns.CreateDeleteDomainRecordRequest()
			req.RecordId = record.Id
			_, e := da.client.DeleteDomainRecord(req)
			if e != nil {
				return e
			}
		}
	}

	return da.AddRecord(dns.AddRecordParams{
		Domain: params.Domain,
		Type:   params.Type,
		Value:  params.Value,
	})
}

func (da *DnsAli) DelRecord(params dns.DelRecordParams) error {
	exists, e := da.QueryRecords(dns.QueryRecordParams{
		Domain: params.Domain,
		Type:   params.Type,
	})
	if e != nil {
		return e
	}
	if len(exists) == 0 {
		return nil
	}
	for _, record := range exists {
		req := alidns.CreateDeleteDomainRecordRequest()
		req.RecordId = record.Id
		_, e := da.client.DeleteDomainRecord(req)
		if e != nil {
			return e
		}
	}
	return nil
}
