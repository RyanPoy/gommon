package ip

import (
	"bufio"
	"encoding/binary"
	"gommon/extends/array"
	"gommon/ip/consts"
	"net"
	"os"
	"sort"
	"strings"
)

func isV4(ipStr string) bool {
	return strings.Contains(ipStr, ".")
}

type IPTable struct {
	v4s       V4Ranges
	v6s       V6Ranges
	countries *array.Array
	isps      *array.Array
	provs     *array.Array
	cities    *array.Array
	numbers   *array.Array
}

func (t *IPTable) AddV4(x *V4Range) {
	t.v4s = append(t.v4s, x)
}

func (t *IPTable) AddV6(x *V6Range) {
	t.v6s = append(t.v6s, x)
}

func (t *IPTable) StringOfV4(ipRange *V4Range) string {
	return ipRange.StartStr + "|" +
		ipRange.EndStr + "|" +
		t.countries.Get(ipRange.CountryIdx) + "|" +
		t.isps.Get(ipRange.IspIdx) + "|" +
		t.provs.Get(ipRange.ProvIdx) + "|" +
		t.cities.Get(ipRange.CityIdx) + "|" +
		t.numbers.Get(ipRange.NumberIdx)
}
func (t *IPTable) StringOfV6(ipRange *V6Range) string {
	return ipRange.StartStr + "|" +
		ipRange.EndStr + "|" +
		t.countries.Get(ipRange.CountryIdx) + "|" +
		t.isps.Get(ipRange.IspIdx) + "|" +
		t.provs.Get(ipRange.ProvIdx) + "|" +
		t.cities.Get(ipRange.CityIdx) + "|" +
		t.numbers.Get(ipRange.NumberIdx)
}
func (t *IPTable) AreaOf(ipRange *IPRange) map[string]string {
	countryCode := t.countries.Get(ipRange.CountryIdx)
	country := consts.CountryOf(countryCode)
	return map[string]string{
		"country":   country.Name,
		"timezone":  country.Timezone,
		"continent": country.ContinentCode,
		"province":  t.provs.Get(ipRange.ProvIdx),
		"city":      t.cities.Get(ipRange.CityIdx),
		"isp":       t.isps.Get(ipRange.IspIdx),
	}
}

func (t *IPTable) SearchV4(ipv4 string) *V4Range {
	ip := net.ParseIP(ipv4)
	if ip == nil {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	ranges := t.v4s
	idx := sort.Search(len(ranges), func(i int) bool {
		return ranges[i].Gte(ip)
		//return cmp(ranges[i].Low, ip) == 1 || cmp(ranges[i].High, ip) != -1
	})

	if idx < len(ranges) && ranges[idx].Contains(ip) {
		return ranges[idx]
	}
	return nil
}

func (t *IPTable) SearchV6(ipv6 string) *V6Range {
	ip := net.ParseIP(ipv6)
	if ip == nil {
		return nil
	}
	ip = ip.To16()
	if ip == nil {
		return nil
	}
	ranges := t.v6s
	idx := sort.Search(len(ranges), func(i int) bool {
		return ranges[i].Gte(ip)
		//return cmp(ranges[i].Low, ip) == 1 || cmp(ranges[i].High, ip) != -1
	})

	if idx < len(ranges) && ranges[idx].Contains(ip) {
		return ranges[idx]
	}
	return nil
}

func (t *IPTable) sortAndUniq() {
	sort.Sort(&t.v4s)
	sort.Sort(&t.v6s)

	uniqV4s, uniqV6s := make(V4Ranges, 0), make(V6Ranges, 0)
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
		v4s:       make(V4Ranges, 0),
		v6s:       make(V6Ranges, 0),
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
			v1, v2 := binary.BigEndian.Uint32(low), binary.BigEndian.Uint32(high)
			table.AddV4(&V4Range{
				Low:  v1,
				High: v2,
				IPRange: IPRange{
					StartStr:   parts[0],
					EndStr:     parts[1],
					CountryIdx: table.countries.Append(parts[2]),
					IspIdx:     table.isps.Append(parts[3]),
					ProvIdx:    table.provs.Append(parts[4]),
					CityIdx:    table.cities.Append(parts[5]),
					NumberIdx:  table.numbers.Append(parts[6]),
				},
			})
		} else {
			table.AddV6(&V6Range{
				Low:  low,
				High: high,
				IPRange: IPRange{
					StartStr:   parts[0],
					EndStr:     parts[1],
					CountryIdx: table.countries.Append(parts[2]),
					IspIdx:     table.isps.Append(parts[3]),
					ProvIdx:    table.provs.Append(parts[4]),
					CityIdx:    table.cities.Append(parts[5]),
					NumberIdx:  table.numbers.Append(parts[6]),
				},
			})
		}
	}
	table.sortAndUniq()
	return table, nil
}
