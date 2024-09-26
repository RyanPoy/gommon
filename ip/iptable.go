package ip

import (
	"gommon/ip/internal"
	"sort"
)

type IPTable struct {
	data      []IPRange
	countries *internal.Array
	isps      *internal.Array
	provs     *internal.Array
	cities    *internal.Array
	numbers   *internal.Array

	searchFunc func(ipStr string, table *IPTable) IPRange
}

func (t *IPTable) Add(x IPRange) {
	t.data = append(t.data, x)
}

func (t *IPTable) Len() int {
	return len(t.data)
}

func (t *IPTable) Swap(i, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
}

func (t *IPTable) Less(i, j int) bool {
	return t.data[i].Cmp(t.data[j]) < 0
}

func (t *IPTable) StringOf(ipRange IPRange) string {
	return ipRange.StartStr() + "|" +
		ipRange.EndStr() + "|" +
		t.countries.Get(ipRange.CountryIdx()) + "|" +
		t.isps.Get(ipRange.IspIdx()) + "|" +
		t.provs.Get(ipRange.ProvIdx()) + "|" +
		t.cities.Get(ipRange.CityIdx()) + "|" +
		t.numbers.Get(ipRange.NumberIdx())
}

func (t *IPTable) Search(ipStr string) IPRange {
	return t.searchFunc(ipStr, t)
}

func NewV4Table(fpath string) (*IPTable, error) {
	table := &IPTable{
		data:      make([]IPRange, 0),
		countries: internal.NewArray(),
		isps:      internal.NewArray(),
		provs:     internal.NewArray(),
		cities:    internal.NewArray(),
		numbers:   internal.NewArray(),
		searchFunc: func(ipStr string, table *IPTable) IPRange {
			ip := uint32Of(ipStr)
			if ip == 0 {
				return nil
			}
			idx := sort.Search(len(table.data), func(i int) bool {
				return table.data[i].GTE(ipStr)
			})

			if idx < len(table.data) && table.data[idx].Contains(&ip) {
				return table.data[idx]
			}
			return nil
		},
	}
	return newTable(fpath, table, ParseV4Range)
}

func NewV6Table(fpath string) (*IPTable, error) {
	table := &IPTable{
		data:      make([]IPRange, 0),
		countries: internal.NewArray(),
		isps:      internal.NewArray(),
		provs:     internal.NewArray(),
		cities:    internal.NewArray(),
		numbers:   internal.NewArray(),
		searchFunc: func(ipStr string, table *IPTable) IPRange {
			ipv := internal.FromIpv6(ipStr)
			if ipv == nil {
				return nil
			}
			idx := sort.Search(len(table.data), func(i int) bool {
				return table.data[i].GTE(ipStr)
			})

			if idx < len(table.data) && table.data[idx].Contains(ipv) {
				return table.data[idx]
			}
			return nil
		},
	}
	return newTable(fpath, table, ParseV6Range)
}

func newTable(fpath string, table *IPTable, parseRange func(string, *IPTable) IPRange) (*IPTable, error) {

	lines, err := internal.LoadFile(fpath)
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		v4Range := parseRange(line, table)
		if v4Range == nil {
			continue
		}
		table.Add(v4Range) // 添加到集合
	}
	sort.Sort(table)
	return table, nil
}
