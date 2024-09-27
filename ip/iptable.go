package ip

import (
	"gommon/extends"
	"sort"
)

type IPRanges struct {
	data []*IPRange
}

func (t *IPRanges) Len() int {
	return len(t.data)
}

func (t *IPRanges) Swap(i, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
}

func (t *IPRanges) Less(i, j int) bool {
	return t.data[i].Cmp(t.data[j]) < 0
}

type IPTable struct {
	data      *IPRanges
	countries *extends.Array
	isps      *extends.Array
	provs     *extends.Array
	cities    *extends.Array
	numbers   *extends.Array

	searchFunc func(ipStr string, ranges *IPRanges) *IPRange
}

func (t *IPTable) Add(x *IPRange) {
	t.data.data = append(t.data.data, x)
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
	return t.searchFunc(ipStr, t.data)
}

func NewV4Table(fpath string) (*IPTable, error) {
	table := &IPTable{
		data:       &IPRanges{data: make([]*IPRange, 0)},
		countries:  extends.NewArray(),
		isps:       extends.NewArray(),
		provs:      extends.NewArray(),
		cities:     extends.NewArray(),
		numbers:    extends.NewArray(),
		searchFunc: SearchV4,
	}
	return newTable(fpath, table, ParseV4Range)
}

func NewV6Table(fpath string) (*IPTable, error) {
	table := &IPTable{
		data:       &IPRanges{data: make([]*IPRange, 0)},
		countries:  extends.NewArray(),
		isps:       extends.NewArray(),
		provs:      extends.NewArray(),
		cities:     extends.NewArray(),
		numbers:    extends.NewArray(),
		searchFunc: SearchV6,
	}
	return newTable(fpath, table, ParseV6Range)
}

func newTable(fpath string, table *IPTable, parseRange func(string, *IPTable) *IPRange) (*IPTable, error) {
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
	sort.Sort(table.data)
	return table, nil
}
