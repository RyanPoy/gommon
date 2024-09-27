package ip

import (
	"bufio"
	"gommon/extends/array"
	"net"
	"os"
	"sort"
	"strings"
)

func isV4(ipStr string) bool {
	return strings.Contains(ipStr, ".")
}

type IPTable struct {
	v4s       IPRanges
	v6s       IPRanges
	countries *array.Array
	isps      *array.Array
	provs     *array.Array
	cities    *array.Array
	numbers   *array.Array
}

func (t *IPTable) AddV4(x *IPRange) {
	t.v4s = append(t.v4s, x)
}

func (t *IPTable) AddV6(x *IPRange) {
	t.v6s = append(t.v6s, x)
}

func (t *IPTable) StringOf(ipRange *IPRange) string {
	return ipRange.StartStr + "|" +
		ipRange.EndStr + "|" +
		t.countries.Get(ipRange.CountryIdx) + "|" +
		t.isps.Get(ipRange.IspIdx) + "|" +
		t.provs.Get(ipRange.ProvIdx) + "|" +
		t.cities.Get(ipRange.CityIdx) + "|" +
		t.numbers.Get(ipRange.NumberIdx)
}

func (t *IPTable) Search(ipStr string) *IPRange {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	if strings.Contains(ipStr, ".") {
		return t.SearchV4(ip)
	} else {
		return t.SearchV6(ip)
	}
}

func (t *IPTable) SearchV4(ipv4 net.IP) *IPRange {
	if ipv4 = ipv4.To4(); ipv4 == nil {
		return nil
	}
	return t.search(ipv4, t.v4s)
}

func (t *IPTable) SearchV6(ipv6 net.IP) *IPRange {
	if ipv6 = ipv6.To16(); ipv6 == nil {
		return nil
	}
	return t.search(ipv6, t.v6s)
}

func (t *IPTable) search(ip net.IP, ranges IPRanges) *IPRange {
	if ip == nil {
		return nil
	}
	idx := sort.Search(len(ranges), func(i int) bool {
		return cmp(ranges[i].low, ip) == 1 || cmp(ranges[i].high, ip) != -1
	})

	if idx < len(ranges) && ranges[idx].Contains(ip) {
		return ranges[idx]
	}
	return nil
}

func (t *IPTable) sortAndUniq() {
	sort.Sort(&t.v4s)
	sort.Sort(&t.v6s)

	uniqV4s, uniqV6s := make(IPRanges, 0), make(IPRanges, 0)
	if len(t.v4s) > 0 {
		uniqV4s = append(uniqV4s, t.v4s[0])
		for i := 1; i < len(t.v4s); i++ {
			if t.v4s[i].Cmp(t.v4s[i-1]) != 0 {
				uniqV4s = append(uniqV4s, t.v4s[i])
			}
		}
	}
	if len(t.v6s) > 0 {
		uniqV6s = append(uniqV6s, t.v6s[0])
		for i := 1; i < len(t.v6s); i++ {
			if t.v6s[i].Cmp(t.v6s[i-1]) != 0 {
				uniqV6s = append(uniqV6s, t.v6s[i])
			}
		}
	}
	t.v4s, t.v6s = uniqV4s, uniqV6s
}

func NewIPTable(fpaths ...string) (*IPTable, error) {
	var err error
	table := &IPTable{
		v4s:       make(IPRanges, 0),
		v6s:       make(IPRanges, 0),
		countries: array.New(),
		isps:      array.New(),
		provs:     array.New(),
		cities:    array.New(),
		numbers:   array.New(),
	}

	for _, fpath := range fpaths {
		if table, err = initFromFile(fpath, table); err != nil {
			return nil, err
		}
	}

	return table, nil
}

func initFromFile(fpath string, table *IPTable) (*IPTable, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	for scanner := bufio.NewScanner(f); scanner.Scan(); {
		line := scanner.Text()
		if line[0] == '#' {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 7 {
			continue
		}

		low, high := net.ParseIP(parts[0]), net.ParseIP(parts[1])
		if low == nil && high == nil {
			continue
		}

		isV4 := isV4(parts[0])
		if isV4 {
			low, high = low.To4(), high.To4()
		} else {
			low, high = low.To16(), high.To16()
		}
		if low == nil || high == nil {
			continue
		}
		if cmp(low, high) == 1 {
			low, high = high, low
		}
		if isV4 {
			table.AddV4(&IPRange{
				low:        low,
				high:       high,
				StartStr:   parts[0],
				EndStr:     parts[1],
				CountryIdx: table.countries.Append(parts[2]),
				IspIdx:     table.isps.Append(parts[3]),
				ProvIdx:    table.provs.Append(parts[4]),
				CityIdx:    table.cities.Append(parts[5]),
				NumberIdx:  table.numbers.Append(parts[6]),
			})
		} else {
			table.AddV6(&IPRange{
				low:        low,
				high:       high,
				StartStr:   parts[0],
				EndStr:     parts[1],
				CountryIdx: table.countries.Append(parts[2]),
				IspIdx:     table.isps.Append(parts[3]),
				ProvIdx:    table.provs.Append(parts[4]),
				CityIdx:    table.cities.Append(parts[5]),
				NumberIdx:  table.numbers.Append(parts[6]),
			})
		}
	}
	table.sortAndUniq()
	return table, nil
}
