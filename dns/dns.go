package dns

const RECORD_TYPE_A = "A"
const RECORD_TYPE_AAAA = "AAAA"
const RECORD_TYPE_CNAME = "CNAME"
const RECORD_TYPE_TXT = "TXT"

type DomainResolved struct {
	// ali: 主域名(veikr.com)，   cloudflare: zoneid
	DomainName string
	RR         string
}

type Record struct {
	Id      string
	Domain  string
	Type    string
	Value   string
	Enabled bool
	Proxied bool
}

type QueryRecordParams struct {
	Domain string
	Type   string
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

type IDns interface {
	ListMainDomains() ([]string, error)
	ListAllRecords(domain string) ([]Record, error)
	QueryRecords(params QueryRecordParams) ([]Record, error)
	AddRecord(params AddRecordParams) error
	EditRecord(params EditRecordParams) error
	DelRecord(params DelRecordParams) error
}
