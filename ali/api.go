package ali

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

type _API struct {
	AccessSecretId  string
	AccessSecretKey string
	Region          string
}

type Request struct {
	Method string
	Action string
	Query  map[string]string
	Body   []byte
	Result any
}

type Response struct {
	HttpCode int
	Headers  http.Header
	Body     []byte
}

type DescribeDomainRecordsReq struct {
	DomainName string
	PageNumber int
	PageSize   int
	Type       string
}

type DescribeDomainRecordsRes struct {
	RequestId     string
	TotalCount    int64
	PageSize      int
	PageNumber    int
	DomainRecords *DomainRecords
}

type DomainRecords struct {
	Record []*RecordItem
}

type RecordItem struct {
	Status          string
	Type            string
	TTL             int
	RecordId        string
	Priority        int
	RR              string
	DomainName      string
	Weight          int
	Value           string
	CreateTimestamp int64
	UpdateTimestamp int64
}

type AddDomainRecordReq struct {
	DomainName string
	RR         string
	Type       string
	Value      string
	TTL        int
	Priority   int
}

type AddDomainRecordRes struct {
	RequestId string
	RecordId  string
}

type DeleteDomainRecordReq struct {
	RecordId string
}

type DeleteDomainRecordRes struct {
	RequestId string
	RecordId  string
}

type UpdateDomainRecordReq struct {
	RecordId string
	RR       string
	Type     string
	Value    string
	TTL      int
	Priority int
}

type UpdateDomainRecordRes struct {
	RequestId string
	RecordId  string
}

type DescribeSubDomainRecordsReq struct {
	SubDomain  string
	PageNumber int
	PageSize   int
	Type       string
}

type DescribeSubDomainRecordsRes struct {
	RequestId     string
	TotalCount    int64
	PageSize      int
	PageNumber    int
	DomainRecords DomainRecords
}

type DescribeDomainsRes struct {
	RequestId  string
	TotalCount int64
	PageSize   int64
	PageNumber int64
	Domains    struct {
		Domain []struct {
			DomainName string
			DomainId   string
		}
	}
}

var api = _API{}

func buildQueryString(m map[string]string) string {
	var keys = make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	var rtn = ""
	for _, key := range keys {
		var value = strings.TrimSpace(m[key])
		rtn += fmt.Sprintf("%s=%s&", url.QueryEscape(key), url.QueryEscape(value))
	}
	if len(rtn) > 0 {
		rtn = rtn[0 : len(rtn)-1]
	}
	return rtn
}

func buildCanonicalRequest(method string, queryString string, canonicalHeaders string, signedHeaders string, bodyHash string) string {
	var rtn = ""
	rtn += fmt.Sprintf("%s\n", method)
	rtn += fmt.Sprintf("%s\n", "/")
	rtn += fmt.Sprintf("%s\n", queryString)
	rtn += fmt.Sprintf("%s\n", canonicalHeaders)
	rtn += fmt.Sprintf("%s\n", signedHeaders)
	rtn += bodyHash
	return rtn
}

func buildCanonicalHeaders(headers map[string]string) string {
	var keys = make([]string, 0)
	for k, _ := range headers {
		keys = append(keys, strings.ToLower(k))
	}
	slices.Sort(keys)

	var rtn = ""
	for _, key := range keys {
		var value = strings.TrimSpace(headers[key])
		var entry = fmt.Sprintf("%s:%s\n", key, value)
		rtn += entry
	}
	return rtn
}

func buildCanonicalSignedHeaders(headers map[string]string) string {
	var keys = make([]string, 0)
	for k, _ := range headers {
		keys = append(keys, strings.ToLower(k))
	}
	slices.Sort(keys)
	return strings.Join(keys, ";")
}

func hash(body []byte) string {
	var bs = sha256.Sum256(body)
	return hex.EncodeToString(bs[:])
}

func buildSignature(accessKeySecret string, canonicalRequest string) string {
	var toSign = "ACS3-HMAC-SHA256" + "\n" + hash([]byte(canonicalRequest))
	var h = hmac.New(sha256.New, []byte(accessKeySecret))
	h.Write([]byte(toSign))
	var bs = h.Sum(nil)
	return hex.EncodeToString(bs)
}

func buildAuth(accessKeyId string, signedHeaders string, signature string) string {
	return fmt.Sprintf("ACS3-HMAC-SHA256 Credential=%s,SignedHeaders=%s,Signature=%s", accessKeyId, signedHeaders, signature)
}

func (a *_API) Send(req *Request) (*Response, error) {

	if req.Body == nil {
		req.Body = []byte{}
	}
	if req.Method == "" {
		req.Method = http.MethodGet
	}

	var queryString = buildQueryString(req.Query)
	var bodyHash = hash(req.Body)
	var u = fmt.Sprintf("https://alidns.%s.aliyuncs.com/?%s", a.Region, queryString)
	log.Println("ali.send", u)
	var headers = map[string]string{}

	var now = time.Now()
	var nonce = strconv.FormatInt(now.UnixMicro(), 16)
	headers["x-acs-action"] = req.Action
	headers["x-acs-version"] = "2015-01-09"
	headers["x-acs-signature-nonce"] = nonce
	headers["x-acs-date"] = now.UTC().Format(time.RFC3339)
	headers["host"] = fmt.Sprintf("alidns.%s.aliyuncs.com", a.Region)
	headers["x-acs-content-sha256"] = bodyHash
	headers["content-type"] = "applicaton/json"
	log.Println("headers", headers)

	var canonicalHeaders = buildCanonicalHeaders(headers)
	var signedHeaders = buildCanonicalSignedHeaders(headers)
	var canonicalRequest = buildCanonicalRequest(req.Method, queryString, canonicalHeaders, signedHeaders, bodyHash)
	var signature = buildSignature(a.AccessSecretKey, canonicalRequest)
	var auth = buildAuth(a.AccessSecretId, signedHeaders, signature)

	log.Println("signing----------")
	log.Println("canonicalHeaders", canonicalHeaders)
	log.Println("signedHeaders", signedHeaders)
	log.Println("canonicalRequest", canonicalRequest)
	log.Println("signature", signature)
	log.Println("auth", auth)

	headers["Authorization"] = auth

	r, e := http.NewRequest(req.Method, u, bytes.NewReader(req.Body))
	if e != nil {
		log.Println(e)
		return nil, e
	}

	for k, v := range headers {
		r.Header.Set(k, v)
	}
	hres, e := http.DefaultClient.Do(r)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	res := &Response{}
	res.HttpCode = hres.StatusCode
	res.Headers = hres.Header
	bs, e := io.ReadAll(hres.Body)
	if e != nil {
		log.Println(e)
		return res, e
	}
	res.Body = bs
	log.Println(hres.StatusCode, hres.Status)
	log.Println(string(bs))
	if hres.StatusCode != 200 {
		return nil, errors.New(string(bs))
	}
	if req.Result != nil {
		e = json.Unmarshal(res.Body, &req.Result)
		if e != nil {
			return nil, e
		}
	}
	return res, nil
}

func (a *_API) ListMainDomains() ([]string, error) {
	domains := &DescribeDomainsRes{}
	_, e := a.Send(&Request{
		Action: "DescribeDomains",
		Query: map[string]string{
			"PageSize": "100",
		},
		Result: domains,
	})
	if e != nil {
		return nil, e
	}
	rtn := make([]string, 0)
	for _, v := range domains.Domains.Domain {
		rtn = append(rtn, v.DomainName)
	}
	return rtn, nil
}

func (a *_API) List(req DescribeDomainRecordsReq) (*DescribeDomainRecordsRes, error) {
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 50
	}
	var rtn = &DescribeDomainRecordsRes{}
	_, e := a.Send(&Request{
		Action: "DescribeDomainRecords",
		Query: map[string]string{
			"DomainName": req.DomainName,
			"PageNumber": strconv.Itoa(req.PageNumber),
			"PageSize":   strconv.Itoa(req.PageSize),
			"Type":       req.Type,
		},
		Result: rtn,
	})
	return rtn, e
}

func (a *_API) ListSub(req DescribeSubDomainRecordsReq) (*DescribeSubDomainRecordsRes, error) {
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 50
	}
	var rtn = &DescribeSubDomainRecordsRes{}
	_, e := a.Send(&Request{
		Action: "DescribeSubDomainRecords",
		Query: map[string]string{
			"SubDomain":  req.SubDomain,
			"PageNumber": strconv.Itoa(req.PageNumber),
			"PageSize":   strconv.Itoa(req.PageSize),
			"Type":       req.Type,
		},
		Result: rtn,
	})
	return rtn, e
}

func (a *_API) Add(req AddDomainRecordReq) (*AddDomainRecordRes, error) {
	if req.TTL <= 0 {
		req.TTL = 600
	}
	if req.Priority == 0 {
		req.Priority = 1
	}
	var rtn = &AddDomainRecordRes{}
	_, e := a.Send(&Request{
		Action: "AddDomainRecord",
		Query: map[string]string{
			"DomainName": req.DomainName,
			"RR":         req.RR,
			"Type":       req.Type,
			"Value":      req.Value,
			"TTL":        strconv.Itoa(req.TTL),
			"Priority":   strconv.Itoa(req.Priority),
		},
		Result: rtn,
	})
	return rtn, e
}

func (a *_API) Update(req UpdateDomainRecordReq) (*UpdateDomainRecordRes, error) {
	if req.TTL <= 0 {
		req.TTL = 600
	}
	var rtn = &UpdateDomainRecordRes{}
	_, e := a.Send(&Request{
		Action: "UpdateDomainRecord",
		Query: map[string]string{
			"RecordId": req.RecordId,
			"RR":       req.RR,
			"Type":     req.Type,
			"Value":    req.Value,
			"TTL":      strconv.Itoa(req.TTL),
			"Priority": strconv.Itoa(req.Priority),
		},
		Result: rtn,
	})
	return rtn, e
}

func (a *_API) Del(req DeleteDomainRecordReq) (*DeleteDomainRecordRes, error) {
	var rtn = &DeleteDomainRecordRes{}
	_, e := a.Send(&Request{
		Action: "DeleteDomainRecord",
		Query: map[string]string{
			"RecordId": req.RecordId,
		},
		Result: rtn,
	})
	return rtn, e
}
