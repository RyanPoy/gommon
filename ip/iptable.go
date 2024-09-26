package ip

import (
	"gommon/extends"
	"sort"
)

type IPTable struct {
	data      []IPRange
	countries *extends.Array
	isps      *extends.Array
	provs     *extends.Array
	cities    *extends.Array
	numbers   *extends.Array

	searcher Searcher
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
	return t.searcher.Search(ipStr, t)
}

func NewV4Table(fpath string) (*IPTable, error) {
	table := &IPTable{
		data:      make([]IPRange, 0),
		countries: extends.NewArray(),
		isps:      extends.NewArray(),
		provs:     extends.NewArray(),
		cities:    extends.NewArray(),
		numbers:   extends.NewArray(),
		searcher:  &V4Searcher{},
	}
	return newTable(fpath, table, ParseV4Range)
}

func NewV6Table(fpath string) (*IPTable, error) {
	table := &IPTable{
		data:      make([]IPRange, 0),
		countries: extends.NewArray(),
		isps:      extends.NewArray(),
		provs:     extends.NewArray(),
		cities:    extends.NewArray(),
		numbers:   extends.NewArray(),
		searcher:  &V6Searcher{},
	}
	return newTable(fpath, table, ParseV6Range)
}

func newTable(fpath string, table *IPTable, parseRange func(string, *IPTable) IPRange) (*IPTable, error) {
	lines, err := LoadFile(fpath)
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
