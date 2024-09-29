package ip

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

func isV4(ipStr string) bool {
	return strings.Contains(ipStr, ".")
}

type Table struct {
	v4s       V4IntervalList
	v6s       V6IntervalList
	countries *Array
	isps      *Array
	provs     *Array
	cities    *Array
	numbers   *Array
}

func (t *Table) AddV4(x *V4Interval) {
	t.v4s = append(t.v4s, x)
}

func (t *Table) AddV6(x *V6Interval) {
	t.v6s = append(t.v6s, x)
}

func (t *Table) StringOf(ipRange *Interval) string {
	return ipRange.StartStr + "|" +
		ipRange.EndStr + "|" +
		t.countries.Get(ipRange.CountryIdx) + "|" +
		t.isps.Get(ipRange.IspIdx) + "|" +
		t.provs.Get(ipRange.ProvIdx) + "|" +
		t.cities.Get(ipRange.CityIdx) + "|" +
		strconv.Itoa(ipRange.Number)
}

func (t *Table) AreaOf(ipRange *Interval) map[string]string {
	countryCode := t.countries.Get(ipRange.CountryIdx)
	country := CountryOf(countryCode)
	return map[string]string{
		"country":   country.Name,
		"timezone":  country.Timezone,
		"continent": country.ContinentCode,
		"province":  t.provs.Get(ipRange.ProvIdx),
		"city":      t.cities.Get(ipRange.CityIdx),
		"isp":       t.isps.Get(ipRange.IspIdx),
	}
}

func (t *Table) Search(ipStr string) *Interval {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}

	if isV4(ipStr) {
		return &t.SearchV4(ip.To4()).Interval
	}
	return &t.SearchV6(ip.To16()).Interval
}

func (t *Table) SearchV4(ip net.IP) *V4Interval {
	if ip == nil {
		return nil
	}
	return t.v4s.Search(ip)
}

func (t *Table) SearchV6(ip net.IP) *V6Interval {
	if ip == nil {
		return nil
	}
	return t.v6s.Search(ip)

}

func (t *Table) sortAndUniq() {
	sort.Sort(&t.v4s)
	sort.Sort(&t.v6s)

	uniqV4s, uniqV6s := make(V4IntervalList, 0), make(V6IntervalList, 0)
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

func NewIPTable(fpaths ...string) (*Table, error) {
	var err error
	table := &Table{
		v4s:       make(V4IntervalList, 0),
		v6s:       make(V6IntervalList, 0),
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

		isV4 := isV4(parts[0])
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
			table.AddV4(&V4Interval{
				Low:  binary.BigEndian.Uint32(low),
				High: binary.BigEndian.Uint32(high),
				Interval: Interval{
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
			table.AddV6(&V6Interval{
				Low:  u128(low),
				High: u128(high),
				Interval: Interval{
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
