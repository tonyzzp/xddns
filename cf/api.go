package cf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const _CF_API = "https://api.cloudflare.com/client/v4"

type _API struct {
	Token string
}

var api = _API{}

type ApiListParams struct {
	Zone    string
	Name    string
	Page    int
	PerPage int
	Type    string
}

type ApiListRes struct {
	Result     []Record
	ResultInfo ResultInfo `json:"result_info"`
}

type ResultInfo struct {
	Count      int
	Page       int
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
}

type Record struct {
	Content string
	Name    string
	Type    string
	ID      string
	TTL     int
	Proxied bool
}

type CreateParams struct {
	Zone    string `json:"-"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
}

type UpdateParams struct {
	Zone     string `json:"-"`
	RecordID string `json:"-"`
	Content  string `json:"content"`
	Name     string `json:"name"`
	Proxied  bool   `json:"proxied"`
	Type     string `json:"type"`
	TTL      int    `json:"ttl"`
}

type Zone struct {
	Id   string
	Name string
}

type ZonesRes struct {
	Success    bool
	ResultInfo ResultInfo `json:"result_info"`
	Result     []Zone
}

func (api *_API) get(p string, params map[string]string, result any) error {
	var requestUrl = fmt.Sprintf("%s%s", _CF_API, p)
	if len(params) > 0 {
		requestUrl = requestUrl + "?"
		for k, v := range params {
			requestUrl = requestUrl + url.QueryEscape(k) + "=" + url.QueryEscape(v) + "&"
		}
	}
	req, e := http.NewRequest(http.MethodGet, requestUrl, nil)
	if e != nil {
		return e
	}
	req.Method = http.MethodGet
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.Token)
	log.Println("cloudflare.get", requestUrl)
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		log.Println(e)
		return e
	}
	data, e := io.ReadAll(resp.Body)
	if e != nil {
		log.Println(e)
		return e
	}

	log.Println(resp.StatusCode, resp.Status)
	log.Println(string(data))
	if resp.StatusCode != 200 {
		return errors.New(string(data))
	}
	e = json.Unmarshal(data, result)
	return e
}

func (api *_API) post(p string, method string, body any, result any) error {
	var url = fmt.Sprintf("%s%s", _CF_API, p)
	var data []byte
	var e error
	if body != nil {
		data, e = json.Marshal(body)
		if e != nil {
			return e
		}
	}
	var ioBody io.Reader
	if data != nil {
		ioBody = bytes.NewReader(data)
	}
	req, e := http.NewRequest(method, url, ioBody)
	if e != nil {
		return e
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.Token)
	log.Println("cloudflare.post", url)
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		log.Println(e)
		return e
	}
	data, e = io.ReadAll(resp.Body)
	if e != nil {
		log.Println(e)
		return e
	}

	log.Println(resp.StatusCode, resp.Status)
	log.Println(string(data))
	if resp.StatusCode != 200 {
		return errors.New(string(data))
	}
	e = json.Unmarshal(data, result)
	return e
}

func (api *_API) ListZones() (*ZonesRes, error) {
	rtn := &ZonesRes{}
	var e = api.get("/zones", nil, rtn)
	return rtn, e
}

func (api *_API) List(params ApiListParams) (*ApiListRes, error) {
	rtn := &ApiListRes{}
	m := make(map[string]string)
	if params.Name != "" {
		m["name"] = params.Name
	}
	if params.Page != 0 {
		m["page"] = strconv.Itoa(params.Page)
	}
	if params.PerPage != 0 {
		m["per_page"] = strconv.Itoa(params.PerPage)
	}
	if params.Type != "" {
		m["type"] = params.Type
	}
	e := api.get(fmt.Sprintf("/zones/%s/dns_records", params.Zone), m, rtn)
	if e != nil {
		return nil, e
	}
	return rtn, nil
}

func (api *_API) Create(params CreateParams) error {
	var u = fmt.Sprintf("/zones/%s/dns_records", params.Zone)
	var m = make(map[string]any)
	e := api.post(u, http.MethodPost, params, &m)
	if e != nil {
		return e
	}
	if m["success"] == true {
		return nil
	} else {
		return fmt.Errorf("%v", m)
	}
}

func (api *_API) Update(params UpdateParams) error {
	log.Println("cf.Update", params)
	var u = fmt.Sprintf("/zones/%s/dns_records/%s", params.Zone, params.RecordID)
	var m = make(map[string]any)
	e := api.post(u, http.MethodPatch, params, &m)
	if e != nil {
		return e
	}
	if m["success"] == true {
		return nil
	} else {
		return fmt.Errorf("%v", m)
	}
}

func (api *_API) Delete(zone string, recordID string) error {
	var u = fmt.Sprintf("/zones/%s/dns_records/%s", zone, recordID)
	var m = make(map[string]any)
	e := api.post(u, http.MethodDelete, nil, &m)
	return e
}
