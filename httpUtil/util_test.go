package httpUtil

import (
	"testing"
)

func TestIPAddrIsLan(t *testing.T) {
	tcs := []struct {
		ip    string
		isLan bool
	}{
		{"100.117.33.108", false},
		{"100.117.33.82", false},
		{"100.117.33.111", false},
		{"10.10.30.1", false},
		{"110.53.216.134", false},
		{"210.53.216.134", false},
	}
	for _, v := range tcs {
		if IPAddrIsLan(v.ip) != v.isLan {
			t.Errorf("IPAddrIsLan(%s) should be %v", v.ip, v.isLan)
		}
	}
}
