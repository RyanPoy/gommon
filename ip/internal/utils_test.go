package internal_test

import (
	"gommon/ip/internal"
	"testing"
)

func TestNormalizeV6Str(t *testing.T) {
	data := [][]string{
		{"2406:0840:f990:0000:0000:0000:0000:0001", "2406:840:f990::1"},
		{"2a13:1801:018a:00cf:0100:0000:0000:0000", "2a13:1801:18a:cf:100::"},
		{"2001:4860:4860:0000:0000:0000:0000:8888", "2001:4860:4860::8888"},
		{"2001:0db8:0000:0000:0000:0000:0000:0001", "2001:db8::1"},
		{"0000:0000:0000:0000:0000:0000:0000:0000", "::"},
		{"0000:0000:0000:0000:0000:0000:0000:0001", "::1"},
		{"2001:0db8:ffff:0000:0123:4567:89ab:cdef", "2001:db8:ffff::123:4567:89ab:cdef"},
		{"1234:5678:9abc:def0:1234:5678:9abc:def0", "1234:5678:9abc:def0:1234:5678:9abc:def0"},
		{"0001:0000:0000:0000:0000:0000:0000:0001", "1::1"},
		{"0000:0000:0000:0000:0000:0000:0001:0002", "::1:2"},
	}
	for i, v := range data {
		expected, short := v[0], v[1]
		if internal.UInt128Of(expected).Cmp(internal.UInt128Of(short)) != 0 {
			t.Fatalf("[%d] -> expected [%s], but [%s]", i, expected, short)
		}
		//full := internal.NormalizeV6(short)
		//if full != expected {
		//	t.Fatalf("[%d] -> expected [%s], but [%s]", i, expected, full)
		//}
	}
}
