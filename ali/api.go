package ali

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
}

type Request struct {
	Method string
	Action string
	Region string
	Query  map[string]string
	Body   []byte
}

type Response struct {
	HttpCode int
	Headers  http.Header
	Body     []byte
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
	var headers = map[string]string{}

	var now = time.Now()
	var nonce = strconv.FormatInt(now.UnixMicro(), 16)
	headers["x-acs-action"] = req.Action
	headers["x-acs-version"] = "2015-01-09"
	headers["x-acs-signature-nonce"] = nonce
	headers["x-acs-date"] = now.UTC().Format(time.RFC3339)
	headers["host"] = fmt.Sprintf("alidns.%s.aliyuncs.com", req.Region)
	headers["x-acs-content-sha256"] = bodyHash
	headers["content-type"] = "applicaton/json"

	var canonicalHeaders = buildCanonicalHeaders(headers)
	var signedHeaders = buildCanonicalSignedHeaders(headers)
	var canonicalRequest = buildCanonicalRequest(req.Method, queryString, canonicalHeaders, signedHeaders, bodyHash)
	var signature = buildSignature(a.AccessSecretKey, canonicalRequest)
	var auth = buildAuth(a.AccessSecretId, signedHeaders, signature)

	headers["Authorization"] = auth

	var u = fmt.Sprintf("https://alidns.%s.aliyuncs.com/?%s", req.Region, queryString)
	log.Println(u)
	r, e := http.NewRequest(req.Method, u, bytes.NewReader(req.Body))
	if e != nil {
		return nil, e
	}

	for k, v := range headers {
		r.Header.Set(k, v)
	}
	hres, e := http.DefaultClient.Do(r)
	if e != nil {
		return nil, e
	}

	res := &Response{}
	res.HttpCode = hres.StatusCode
	res.Headers = hres.Header
	bs, e := io.ReadAll(hres.Body)
	if e != nil {
		return res, e
	}
	res.Body = bs
	return res, nil
}
