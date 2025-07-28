package cf

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/tonyzzp/xddns/config"
	"github.com/tonyzzp/xddns/dns"
)

type DnsCloudFlare struct {
	ctx     context.Context
	cfg     config.ConfigCloudFlare
	domains []dns.Domain
}

func New() *DnsCloudFlare {
	api.Token = config.Config.CloudFlare.Token
	return &DnsCloudFlare{
		cfg: config.Config.CloudFlare,
		ctx: context.Background(),
	}
}

func (cf *DnsCloudFlare) resolve(domain string) (*dns.DomainResolved, error) {
	all, e := cf.ListMainDomains()
	if e != nil {
		return nil, e
	}
	for _, v := range all {
		if strings.HasSuffix(domain, v.Name) {
			rr := strings.TrimSuffix(domain, v.Name)
			if strings.HasSuffix(rr, ".") {
				rr = rr[0 : len(rr)-1]
			} else if rr == "" {
				rr = "@"
			}
			info := &dns.DomainResolved{
				DomainName: v.Id,
				RR:         rr,
			}
			return info, nil

		}
	}
	return nil, errors.New("domain not exists")
}

func (cf *DnsCloudFlare) ListMainDomains() ([]dns.Domain, error) {
	if cf.domains != nil {
		return cf.domains, nil
	}
	rtn := make([]dns.Domain, 0)
	for {
		res, e := api.ListZones()
		if e != nil {
			return nil, e
		}
		for _, v := range res.Result {
			rtn = append(rtn, dns.Domain{
				Name: v.Name,
				Id:   v.Id,
			})
		}
		if len(rtn) >= res.ResultInfo.Count {
			break
		}
	}
	cf.domains = rtn
	return rtn, nil
}

func (cf *DnsCloudFlare) ListAllRecords(domain string) ([]dns.Record, error) {
	log.Println("cf.ListAllRecords", domain)
	info, e := cf.resolve(domain)
	if e != nil {
		return nil, e
	}
	res, e := api.List(ApiListParams{Zone: info.DomainName})
	if e != nil {
		return nil, e
	}
	rtn := make([]dns.Record, 0)
	for _, v := range res.Result {
		rtn = append(rtn, dns.Record{
			Id:      v.ID,
			Domain:  v.Name,
			Type:    string(v.Type),
			Value:   v.Content,
			Enabled: true,
			Proxied: v.Proxied,
			TLT:     v.TTL,
		})
	}
	return rtn, nil
}

func (cf *DnsCloudFlare) QueryRecords(params dns.QueryRecordParams) ([]dns.Record, error) {
	log.Println("cf.QueryRecords", params)
	info, e := cf.resolve(params.Domain)
	if e != nil {
		return nil, e
	}
	res, e := api.List(ApiListParams{
		Zone: info.DomainName,
		Type: params.Type,
		Name: params.Domain,
	})
	if e != nil {
		return nil, e
	}
	rtn := make([]dns.Record, 0)
	for _, v := range res.Result {
		rtn = append(rtn, dns.Record{
			Id:      v.ID,
			Domain:  v.Name,
			Type:    string(v.Type),
			Value:   v.Content,
			Enabled: true,
			Proxied: v.Proxied,
		})
	}
	return rtn, nil
}

func (cf *DnsCloudFlare) AddRecord(params dns.AddRecordParams) error {
	log.Println("cf.AddRecord", params)
	info, e := cf.resolve(params.Domain)
	if e != nil {
		return e
	}
	args := CreateParams{
		Zone:    info.DomainName,
		Content: params.Value,
		Name:    params.Domain,
		Type:    params.Type,
	}
	if params.TTL > 0 {
		args.TTL = params.TTL
	}
	e = api.Create(args)
	return e
}

func (cf *DnsCloudFlare) EditRecord(params dns.EditRecordParams) error {
	log.Println("cf.EditRecord", params)
	info, e := cf.resolve(params.Domain)
	if e != nil {
		log.Println(e)
		return e
	}
	exist, e := cf.QueryRecords(dns.QueryRecordParams{
		Domain: params.Domain,
		Type:   params.Type,
	})
	if e != nil {
		log.Println(e)
		return e
	}
	if len(exist) == 0 {
		args := dns.AddRecordParams{
			Domain: params.Domain,
			Type:   params.Type,
			Value:  params.Value,
			TTL:    params.TTL,
		}
		if params.TTL > 0 {
			args.TTL = params.TTL
		}
		return cf.AddRecord(args)
	} else {
		first := exist[0]
		args := UpdateParams{
			Zone:     info.DomainName,
			RecordID: first.Id,
			Type:     params.Type,
			Content:  params.Value,
			TTL:      params.TTL,
		}
		if params.TTL > 0 {
			args.TTL = params.TTL
		}
		e = api.Update(args)
		return e
	}
}

func (cf *DnsCloudFlare) DelRecord(params dns.DelRecordParams) error {
	log.Println("cf.DelRecord", params)
	exist, e := cf.QueryRecords(dns.QueryRecordParams{
		Domain: params.Domain,
		Type:   params.Type,
	})
	if e != nil {
		return e
	}
	if len(exist) == 0 {
		return nil
	}
	for _, v := range exist {
		info, e := cf.resolve(v.Domain)
		if e != nil {
			return e
		}
		e = api.Delete(info.DomainName, v.Id)
		if e != nil {
			return e
		}
	}
	return nil
}
