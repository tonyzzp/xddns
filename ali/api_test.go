package ali

import (
	"testing"
)

func TestAPI(t *testing.T) {
	api.AccessSecretId = ""
	api.AccessSecretKey = ""

	res, e := api.Send(&Request{
		Region: "cn-shenzhen",
		Action: "DescribeDomainRecords",
		Query: map[string]string{
			"DomainName": "veikr.com",
		},
	})
	t.Log(e)
	if res != nil {
		t.Log(res.HttpCode)
		t.Log(res.Headers)
		t.Log(string(res.Body))
	}
	t.FailNow()
}
