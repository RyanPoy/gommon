package ip

import (
	"gommon/extends"
	"net"
	"sort"
	"strings"
)

func ipv6(ipStr string) net.IP {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	return ip.To16()
}

func ipv4(ipStr string) net.IP {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	return ip.To4()
}

func isV4(ipStr string) bool {
	return strings.Contains(ipStr, ".")
}

type IPTable struct {
	v4s       IPRanges
	v6s       IPRanges
	countries *extends.Array
	isps      *extends.Array
	provs     *extends.Array
	cities    *extends.Array
	numbers   *extends.Array
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
	if strings.Contains(ipStr, ".") {
		ip := ipv4(ipStr)
		return t.SearchV4(ip)
	} else {
		ip := ipv6(ipStr)
		return t.SearchV6(ip)
	}
}

func (t *IPTable) SearchV4(ip net.IP) *IPRange {
	return t.search(ip, t.v4s)
}

func (t *IPTable) SearchV6(ip net.IP) *IPRange {
	return t.search(ip, t.v6s)
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

func NewTable(fpath string) (*IPTable, error) {
	lines, err := LoadFile(fpath)
	if err != nil {
		return nil, err
	}

	table := &IPTable{
		v4s:       make(IPRanges, 0),
		v6s:       make(IPRanges, 0),
		countries: extends.NewArray(),
		isps:      extends.NewArray(),
		provs:     extends.NewArray(),
		cities:    extends.NewArray(),
		numbers:   extends.NewArray(),
	}

	for _, line := range lines {
		if line[0] == '#' {
			continue
		}
		vs := strings.Split(line, "|")
		if len(vs) != 7 {
			continue
		}
		var low, high net.IP
		isV4 := isV4(vs[0])
		if isV4 {
			low, high = ipv4(vs[0]), ipv4(vs[1])
		} else {
			low, high = ipv6(vs[0]), ipv6(vs[1])
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
				StartStr:   vs[0],
				EndStr:     vs[1],
				CountryIdx: table.countries.Append(vs[2]),
				IspIdx:     table.isps.Append(vs[3]),
				ProvIdx:    table.provs.Append(vs[4]),
				CityIdx:    table.cities.Append(vs[5]),
				NumberIdx:  table.numbers.Append(vs[6]),
			})
		} else {
			table.AddV6(&IPRange{
				low:        low,
				high:       high,
				StartStr:   vs[0],
				EndStr:     vs[1],
				CountryIdx: table.countries.Append(vs[2]),
				IspIdx:     table.isps.Append(vs[3]),
				ProvIdx:    table.provs.Append(vs[4]),
				CityIdx:    table.cities.Append(vs[5]),
				NumberIdx:  table.numbers.Append(vs[6]),
			})
		}
	}
	sort.Sort(&table.v4s)
	sort.Sort(&table.v6s)
	return table, nil
}
