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
	ctx context.Context
	cfg config.ConfigCloudFlare
}

func New() *DnsCloudFlare {
	api.Token = config.Config.CloudFlare.Token
	return &DnsCloudFlare{
		cfg: config.Config.CloudFlare,
		ctx: context.Background(),
	}
}

func (cf *DnsCloudFlare) resolve(domain string) (*dns.DomainResolved, error) {
	for v, zone := range cf.cfg.Domains {
		if strings.HasSuffix(domain, v) {
			rr := strings.TrimSuffix(domain, v)
			if strings.HasSuffix(rr, ".") {
				rr = rr[0 : len(rr)-1]
			} else if rr == "" {
				rr = "@"
			}
			info := &dns.DomainResolved{
				DomainName: zone,
				RR:         rr,
			}
			return info, nil
		}
	}
	return nil, errors.New("not exists")
}

func (cf *DnsCloudFlare) ListMainDomains() ([]string, error) {
	rtn := make([]string, 0)
	for k := range cf.cfg.Domains {
		rtn = append(rtn, k)
	}
	return rtn, nil
}

func (cf *DnsCloudFlare) ListAllRecords(domain string) ([]dns.Record, error) {
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
		})
	}
	return rtn, nil
}

func (cf *DnsCloudFlare) QueryRecords(params dns.QueryRecordParams) ([]dns.Record, error) {
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
	e = api.Create(CreateParams{
		Zone:    info.DomainName,
		Content: params.Value,
		Name:    params.Domain,
		Type:    params.Type,
	})
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
		return cf.AddRecord(dns.AddRecordParams{
			Domain: params.Domain,
			Type:   params.Type,
			Value:  params.Value,
		})
	} else {
		first := exist[0]
		e = api.Update(UpdateParams{
			Zone:     info.DomainName,
			RecordID: first.Id,
			Type:     params.Type,
			Content:  params.Value,
		})
		return e
	}
}

func (cf *DnsCloudFlare) DelRecord(params dns.DelRecordParams) error {
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
