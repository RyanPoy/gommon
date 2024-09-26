package ip_test

import (
	"gommon/ip"
	"gommon/ip/internal"
	"strings"
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

	expected := "240e:6af:4700::|240e:6af:47ff:ffff:ffff:ffff:ffff:ffff|CN|CT|江苏|淮安|95566"
	if v6s.StringOf(v6) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, v6s.StringOf(v6))
		return
	}
}

func TestV6ComplexSearch(t *testing.T) {
	fpath := "./test_data/mgiplib-v6-std.txt.latest"
	v6s, err := ip.NewV6s(fpath)
	if err != nil {
		t.Errorf("Can not load v6 data file, %v", err)
		return
	}
	searchIps, err := getAllIpv6(fpath)
	if err != nil {
		t.Errorf("Can not read ipv6 data, %v", err)
		return
	}

	for _, ip := range searchIps {
		v6 := v6s.Search(ip)
		if v6 == nil {
			t.Errorf("Can not find [%s]", "240e:6af:4700:1111::")
			return
		}
		if v6.StartStr != ip {
			t.Errorf("Expected[%s], but[%s]", ip, v6.StartStr)
			return
		}
	}
}

func getAllIpv6(fpath string) ([]string, error) {
	lines, err := internal.LoadFile(fpath)
	if err != nil {
		return nil, err
	}
	relt := make([]string, 0)
	for _, l := range lines {
		if len(l) > 0 && l[0] == '#' {
			continue
		}
		parts := strings.Split(l, "|")
		relt = append(relt, parts[0])
	}
	return relt, nil
}
