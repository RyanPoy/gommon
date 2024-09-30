package ip

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"strconv"
	"strings"
)

func IsV4(ipStr string) bool {
	return strings.Contains(ipStr, ".")
}

type Table struct {
	v4s       IntervalList
	v6s       IntervalList
	countries *Array
	isps      *Array
	provs     *Array
	cities    *Array
	numbers   *Array
}

func (t *Table) StringOf(ipRange Interval) string {
	base := ipRange.BaseInfo()
	return base.StartStr + "|" +
		base.EndStr + "|" +
		t.countries.Get(base.CountryIdx) + "|" +
		t.isps.Get(base.IspIdx) + "|" +
		t.provs.Get(base.ProvIdx) + "|" +
		t.cities.Get(base.CityIdx) + "|" +
		strconv.Itoa(base.Number)
}

func (t *Table) AreaOf(ipRange Interval) map[string]string {
	base := ipRange.BaseInfo()
	countryCode := t.countries.Get(base.CountryIdx)
	country := CountryOf(countryCode)
	return map[string]string{
		"country":   country.Name,
		"timezone":  country.Timezone,
		"continent": country.ContinentCode,
		"province":  t.provs.Get(base.ProvIdx),
		"city":      t.cities.Get(base.CityIdx),
		"isp":       t.isps.Get(base.IspIdx),
	}
}

func (t *Table) Search(ipStr string) Interval {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}

	isV4 := IsV4(ipStr)
	if isV4 {
		if ip = ip.To4(); ip == nil {
			return nil
		}
		return t.v4s.Search(ip)
	} else {
		if ip = ip.To16(); ip == nil {
			return nil
		}
		return t.v6s.Search(ip)
	}
}

func (t *Table) SearchV4(ip net.IP) Interval {
	if ip == nil {
		return nil
	}
	return t.v4s.Search(ip)
}

func (t *Table) SearchV6(ip net.IP) Interval {
	if ip == nil {
		return nil
	}
	return t.v6s.Search(ip)

}

func (t *Table) sortAndUniq() {
	t.v4s.Sort()
	t.v6s.Sort()
	uniqV4s, uniqV6s := make(IntervalList, 0), make(IntervalList, 0)
	if len(t.v4s) > 0 {
		uniqV4s.Add(t.v4s[0])
		for i := 1; i < len(t.v4s); i++ {
			if t.v4s[i].Cmp(t.v4s[i-1]) != 0 {
				uniqV4s.Add(t.v4s[i])
			}
		}
	}
	if len(t.v6s) > 0 {
		uniqV6s.Add(t.v6s[0])
		for i := 1; i < len(t.v6s); i++ {
			if t.v6s[i].Cmp(t.v6s[i-1]) != 0 {
				uniqV6s.Add(t.v6s[i])
			}
		}
	}
	t.v4s, t.v6s = uniqV4s, uniqV6s
}

func NewIPTable(fpaths ...string) (*Table, error) {
	var err error
	table := &Table{
		v4s:       make(IntervalList, 0),
		v6s:       make(IntervalList, 0),
		countries: array(),
		isps:      array(),
		provs:     array(),
		cities:    array(),
		numbers:   array(),
	}

	for _, fpath := range fpaths {
		if table, err = initFromFile(fpath, table); err != nil {
			return nil, err
		}
	}

	return table, nil
}

func initFromFile(fpath string, table *Table) (*Table, error) {
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

		isV4 := IsV4(parts[0])
		if isV4 {
			low, high = low.To4(), high.To4()
		} else {
			low, high = low.To16(), high.To16()
		}
		if low == nil || high == nil {
			continue
		}
		if bytes.Compare(low, high) == 1 {
			low, high = high, low
		}
		number := 0
		number, err = strconv.Atoi(parts[6])
		if err != nil {
			number = 0
		}

		if isV4 {
			table.v4s.Add(&V4Interval{
				Low:  binary.BigEndian.Uint32(low),
				High: binary.BigEndian.Uint32(high),
				baseInfo: &BaseInfo{
					StartStr:   parts[0],
					EndStr:     parts[1],
					CountryIdx: table.countries.Append(parts[2]),
					IspIdx:     table.isps.Append(parts[3]),
					ProvIdx:    table.provs.Append(parts[4]),
					CityIdx:    table.cities.Append(parts[5]),
					Number:     number,
				},
			})
		} else {
			table.v6s.Add(&V6Interval{
				Low:  u128(low),
				High: u128(high),
				baseInfo: &BaseInfo{
					StartStr:   parts[0],
					EndStr:     parts[1],
					CountryIdx: table.countries.Append(parts[2]),
					IspIdx:     table.isps.Append(parts[3]),
					ProvIdx:    table.provs.Append(parts[4]),
					CityIdx:    table.cities.Append(parts[5]),
					Number:     number,
				},
			})
		}
	}
	table.sortAndUniq()
	return table, nil
}
