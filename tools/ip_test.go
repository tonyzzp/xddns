package tools

import (
	"fmt"
	"testing"
)

func TestIp(t *testing.T) {
	fmt.Println(GetExternalIpv4())
	fmt.Println(GetExternalIpv6())
	t.FailNow()
}
