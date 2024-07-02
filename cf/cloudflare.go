package cf

import (
	"ali-ddns/config"
	"ali-ddns/dns"
	"context"
	"errors"
	"strings"

	"github.com/cloudflare/cloudflare-go/v2"
	cfdns "github.com/cloudflare/cloudflare-go/v2/dns"
	"github.com/cloudflare/cloudflare-go/v2/option"
)

type DnsCloudFlare struct {
	ctx    context.Context
	client *cloudflare.Client
	cfg    config.ConfigCloudFlare
}

func New() *DnsCloudFlare {
	return &DnsCloudFlare{
		client: cloudflare.NewClient(
			option.WithAPIToken(config.Config.CloudFlare.Token),
		),
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
	res, e := cf.client.DNS.Records.List(cf.ctx, cfdns.RecordListParams{
		ZoneID: cloudflare.String(info.DomainName),
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
			Value:   v.Content.(string),
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
	res, e := cf.client.DNS.Records.List(cf.ctx, cfdns.RecordListParams{
		ZoneID: cloudflare.String(info.DomainName),
		Type:   cloudflare.F(cfdns.RecordListParamsType(params.Type)),
		Name:   cloudflare.String(params.Domain),
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
			Value:   v.Content.(string),
			Enabled: true,
		})
	}
	return rtn, nil
}

func (cf *DnsCloudFlare) AddRecord(params dns.AddRecordParams) error {
	info, e := cf.resolve(params.Domain)
	if e != nil {
		return e
	}
	_, e = cf.client.DNS.Records.New(cf.ctx, cfdns.RecordNewParams{
		ZoneID: cloudflare.String(info.DomainName),
		Record: cfdns.RecordParam{
			Name:    cloudflare.String(info.RR),
			Type:    cloudflare.F(cfdns.RecordType(params.Type)),
			Content: cloudflare.Raw[any](params.Value),
		},
	})
	return e
}

func (cf *DnsCloudFlare) EditRecord(params dns.EditRecordParams) error {
	info, e := cf.resolve(params.Domain)
	if e != nil {
		return e
	}
	exist, e := cf.QueryRecords(dns.QueryRecordParams{
		Domain: params.Domain,
		Type:   params.Type,
	})
	if e != nil {
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
		_, e = cf.client.DNS.Records.Update(cf.ctx, first.Id, cfdns.RecordUpdateParams{
			ZoneID: cloudflare.String(info.DomainName),
			Record: cfdns.RecordParam{
				Name:    cloudflare.String(info.RR),
				Type:    cloudflare.F(cfdns.RecordType(params.Type)),
				Content: cloudflare.Raw[any](params.Value),
			},
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
		_, e = cf.client.DNS.Records.Delete(cf.ctx, v.Id, cfdns.RecordDeleteParams{
			ZoneID: cloudflare.String(info.DomainName),
		})
		if e != nil {
			return e
		}
	}
	return nil
}
