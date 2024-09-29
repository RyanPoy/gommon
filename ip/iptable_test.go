package ip_test

import (
	"bufio"
	"gommon/ip"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	simple_v4_fpath = "./test_data/a.v4.txt"
	full_v4_fpath   = "./test_data/mgiplib-std.txt.latest"
	full_v6_fpath   = "./test_data/mgiplib-v6-std.txt.latest"
)

var fpaths = []string{simple_v4_fpath, full_v4_fpath, full_v6_fpath}

func TestInitTable(t *testing.T) {
	if _, err := ip.NewIPTable(fpaths...); err != nil {
		t.Errorf("Can not load file, %v", err)
	}

	if _, err := loadIP(fpaths...); err != nil {
		t.Errorf("Can not read ips, %v", err)
	}
}
func TestSearch(t *testing.T) {
	table, _ := ip.NewIPTable(fpaths...)
	for _, caze := range [][]interface{}{
		{
			"0.0.0.0",
			"0.0.0.0|0.255.255.255|HW|OTHER|海外|未知|1",
			map[string]string{"city": "未知", "continent": "", "country": "HW", "isp": "OTHER", "province": "海外", "timezone": ""},
		},
		{
			"223.242.47.30",
			"223.242.32.0|223.242.47.255|CN|CT|安徽|芜湖|576074",
			map[string]string{"city": "芜湖", "continent": "AS", "country": "中国", "isp": "CT", "province": "安徽", "timezone": "Asia/Shanghai"},
		},
		{
			"240e:6af:4700:1111::",
			"240e:6af:4700::|240e:6af:47ff:ffff:ffff:ffff:ffff:ffff|CN|CT|江苏|淮安|95566",
			map[string]string{"city": "淮安", "continent": "AS", "country": "中国", "isp": "CT", "province": "江苏", "timezone": "Asia/Shanghai"},
		},
	} {
		ipStr, expected, expectedMap := caze[0].(string), caze[1].(string), caze[2].(map[string]string)
		ipRange := table.Search(ipStr)
		if ipRange == nil {
			t.Errorf("Expected to get [%s], but not found by [%s]", expected, ipStr)
			return
		}
		if table.StringOf(ipRange) != expected {
			t.Errorf("Expected to get [%s], but got [%s]", expected, table.StringOf(ipRange))
			return
		}
		if !reflect.DeepEqual(table.AreaOf(ipRange), expectedMap) {
			t.Errorf("Expected to get [%v], but got [%v]", expectedMap, table.AreaOf(ipRange))
			return
		}
	}
}

func TestSearchMissing(t *testing.T) {
	table, _ := ip.NewIPTable(fpaths...)
	for _, missing := range []string{"223.242.64.289"} {
		if ipRange := table.Search(missing); ipRange != nil {
			t.Errorf("Expected not to find, but got [%s]", table.StringOf(ipRange))
			return
		}
	}
}

func TestSearchAll(t *testing.T) {
	table, _ := ip.NewIPTable(fpaths...)
	searchIps, _ := loadIP(fpaths...)
	for _, ipStr := range searchIps {
		ipRange := table.Search(ipStr)
		if ipRange == nil {
			t.Errorf("Expected to find [%s], but not found", ipStr)
			return
		}
		if ipRange.StartStr != ipStr {
			t.Errorf("Expected to find [%s], but got [%s]", ipStr, ipRange.StartStr)
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

		func(file *os.File) { defer file.Close() }(f)

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

func BenchmarkSearchV4(b *testing.B) {
	table, _ := ip.NewIPTable(fpaths...)
	//ips, _ := loadIP(full_v4_fpath)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			table.Search("27.190.250.164")
			//for _, ipStr := range ips {
			//	table.Search(ipStr)
			//}
		}
	})
}

func BenchmarkSearchV6(b *testing.B) {
	table, _ := ip.NewIPTable(fpaths...)
	//ips, _ := loadIP(full_v6_fpath)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			table.Search("2001:218:0:2000::147")
			//for _, ipStr := range ips {
			//	table.Search(ipStr)
			//}
		}
	})
}
