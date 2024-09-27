package ip_test

import (
	"bufio"
	"gommon/ip"
	"os"
	"strings"
	"testing"
)

func TestV4TableSimpleSearch(t *testing.T) {
	table, err := ip.NewIPTable("./test_data/a.v4.txt")
	if err != nil {
		t.Errorf("Can not load ipRange v4s file, %v", err)
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
	table, err := ip.NewIPTable("./test_data/mgiplib-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load ipRange v4s file, %v", err)
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
	table, err := ip.NewIPTable(fpath)
	if err != nil {
		t.Errorf("Can not load ipRange v4s file, %v", err)
		return
	}
	searchIps, err := loadIP(fpath)
	if err != nil {
		t.Errorf("Can not read ipv6 v4s, %v", err)
		return
	}

	for _, ip := range searchIps {
		ipRange := table.Search(ip)
		if ipRange == nil {
			t.Errorf("Can not find [%s]", ip)
			return
		}
		if ipRange.StartStr != ip {
			t.Errorf("Expected[%s], but[%s]", ip, ipRange.StartStr)
			return
		}
	}
}

func TestV6TableSearch(t *testing.T) {
	table, err := ip.NewIPTable("./test_data/mgiplib-v6-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load v6 v4s file, %v", err)
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
	table, err := ip.NewIPTable(fpath)
	if err != nil {
		t.Errorf("Can not load v6 v4s file, %v", err)
		return
	}
	searchIps, err := loadIP(fpath)
	if err != nil {
		t.Errorf("Can not read ipv6 v4s, %v", err)
		return
	}

	for _, ip := range searchIps {
		ipRange := table.Search(ip)
		if ipRange == nil {
			t.Errorf("Can not find [%s]", "240e:6af:4700:1111::")
			return
		}
		if ipRange.StartStr != ip {
			t.Errorf("Expected[%s], but[%s]", ip, ipRange.StartStr)
			return
		}
	}
}
func loadIP(fpaths ...string) ([]string, error) {
	lines := make([]string, 0)

	for _, fpath := range fpaths {
		f, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		for scanner := bufio.NewScanner(f); scanner.Scan(); {
			line := scanner.Text()
			if line[0] == '#' || len(line) == 0 || !strings.Contains(line, "|") {
				continue
			}
			lines = append(lines, strings.Split(line, "|")[0])
		}
	}

	return lines, nil
}
