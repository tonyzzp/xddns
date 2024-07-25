package ali

import (
	"errors"
	"strings"

	"github.com/tonyzzp/xddns/config"
	"github.com/tonyzzp/xddns/dns"
)

type DnsAli struct {
	cfg config.ConfigAli
}

func New() *DnsAli {
	rtn := &DnsAli{
		cfg: config.Config.Ali,
	}
	api.AccessSecretId = rtn.cfg.KeyId
	api.AccessSecretKey = rtn.cfg.KeySecret
	api.Region = rtn.cfg.Region
	return rtn
}

func (da *DnsAli) resolve(domain string) (*dns.DomainResolved, error) {
	all, e := da.ListMainDomains()
	if e != nil {
		return nil, e
	}
	for _, d := range all {
		if strings.HasSuffix(domain, d.Name) {
			rr := strings.TrimSuffix(domain, d.Name)
			if strings.HasSuffix(rr, ".") {
				rr = rr[0 : len(rr)-1]
			} else if rr == "" {
				rr = "@"
			}
			return &dns.DomainResolved{
				DomainName: d.Name,
				RR:         rr,
			}, nil
		}
	}
	return nil, errors.New("not exist")
}

func (da *DnsAli) ListMainDomains() ([]dns.Domain, error) {
	all, e := api.ListMainDomains()
	if e != nil {
		return nil, e
	}
	rtn := make([]dns.Domain, 0)
	for _, v := range all {
		rtn = append(rtn, dns.Domain{
			Name: v,
			Id:   "",
		})
	}
	return rtn, nil
}

func (da *DnsAli) ListAllRecords(mainDomain string) ([]dns.Record, error) {
	info, e := da.resolve(mainDomain)
	if e != nil {
		return nil, e
	}
	var page = 1
	doReq := func(page int) (int64, []dns.Record, error) {
		res, e := api.List(DescribeDomainRecordsReq{
			DomainName: info.DomainName,
			PageNumber: page,
			PageSize:   50,
		})
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
		if total == int64(len(rtn)) || len(list) == 0 {
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
		res, e := api.ListSub(DescribeSubDomainRecordsReq{
			SubDomain:  params.Domain,
			PageNumber: page,
			PageSize:   50,
			Type:       params.Type,
		})
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
	_, e = api.Add(AddDomainRecordReq{
		DomainName: info.DomainName,
		Type:       params.Type,
		RR:         info.RR,
		Value:      params.Value,
	})
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
			_, e := api.Del(DeleteDomainRecordReq{RecordId: record.Id})
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
		_, e := api.Del(DeleteDomainRecordReq{RecordId: record.Id})
		if e != nil {
			return e
		}
	}
	return nil
}
