package ip_test

import (
	"gommon/ip"
	"strings"
	"testing"
)

func TestV4TableSimpleSearch(t *testing.T) {
	table, err := ip.NewV4Table("./test_data/a.v4.txt")
	if err != nil {
		t.Errorf("Can not load ipRange data file, %v", err)
		return
	}

	ipRange := table.Search("0.0.0.0")
	if ipRange == nil {
		t.Errorf("Can not find [%s]", "0.0.0.0")
		return
	}

	expected := "0.0.0.0|0.255.255.255|HW|OTHER|海外|未知|1"
	if table.StringOf(ipRange) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, table.StringOf(ipRange))
		return
	}
}

func TestV4TableSearch(t *testing.T) {
	table, err := ip.NewV4Table("./test_data/mgiplib-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load ipRange data file, %v", err)
		return
	}

	ipRange := table.Search("223.242.47.30")
	if ipRange == nil {
		t.Errorf("Can not find [%s]", "223.242.32.30")
		return
	}

	expected := "223.242.32.0|223.242.47.255|CN|CT|安徽|芜湖|576074"
	if table.StringOf(ipRange) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, table.StringOf(ipRange))
		return
	}

	if ipRange := table.Search("223.242.64.289"); ipRange != nil {
		t.Errorf("Expected not found, but got [%s]", table.StringOf(ipRange))
		return
	}
}

func TestV4TableComplexSearch(t *testing.T) {
	fpath := "./test_data/mgiplib-std.txt.latest"
	table, err := ip.NewV4Table(fpath)
	if err != nil {
		t.Errorf("Can not load ipRange data file, %v", err)
		return
	}
	searchIps, err := getAllIp(fpath)
	if err != nil {
		t.Errorf("Can not read ipv6 data, %v", err)
		return
	}

	for _, ip := range searchIps {
		ipRange := table.Search(ip)
		if ipRange == nil {
			t.Errorf("Can not find [%s]", ip)
			return
		}
		if ipRange.OriginData().StartStr != ip {
			t.Errorf("Expected[%s], but[%s]", ip, ipRange.OriginData().StartStr)
			return
		}
	}
}

func TestV6TableSearch(t *testing.T) {
	table, err := ip.NewV6Table("./test_data/mgiplib-v6-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load v6 data file, %v", err)
		return
	}

	ipRange := table.Search("240e:6af:4700:1111::")
	if ipRange == nil {
		t.Errorf("Can not find [%s]", "240e:6af:4700:1111::")
		return
	}

	expected := "240e:6af:4700::|240e:6af:47ff:ffff:ffff:ffff:ffff:ffff|CN|CT|江苏|淮安|95566"
	if table.StringOf(ipRange) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, table.StringOf(ipRange))
		return
	}
}

func TestV6TableComplexSearch(t *testing.T) {
	fpath := "./test_data/mgiplib-v6-std.txt.latest"
	table, err := ip.NewV6Table(fpath)
	if err != nil {
		t.Errorf("Can not load v6 data file, %v", err)
		return
	}
	searchIps, err := getAllIp(fpath)
	if err != nil {
		t.Errorf("Can not read ipv6 data, %v", err)
		return
	}

	for _, ip := range searchIps {
		ipRange := table.Search(ip)
		if ipRange == nil {
			t.Errorf("Can not find [%s]", "240e:6af:4700:1111::")
			return
		}
		if ipRange.OriginData().StartStr != ip {
			t.Errorf("Expected[%s], but[%s]", ip, ipRange.OriginData().StartStr)
			return
		}
	}
}

func getAllIp(fpath string) ([]string, error) {
	lines, err := ip.LoadFile(fpath)
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
