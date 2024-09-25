package ip_test

import (
	"gommon/ip"
	"testing"
)

func TestV6Search(t *testing.T) {
	v6s, err := ip.NewV6s("./test_data/mgiplib-v6-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load v6 data file, %v", err)
		return
	}

	v6 := v6s.Search("240e:6af:4700:1111::")
	if v6 == nil {
		t.Errorf("Can not find [%s]", "240e:6af:4700:1111::")
		return
	}

	expected := "240e:06af:4700:0000:0000:0000:0000:0000|240e:6af:47ff:ffff:ffff:ffff:ffff:ffff|CN|CT|江苏|淮安|95566"
	if v6s.StringOf(v6) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, v6s.StringOf(v6))
		return
	}

}
