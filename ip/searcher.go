package ip

import (
	"gommon/ip/internal"
	"sort"
)

type Searcher interface {
	Search(ipStr string, table *IPTable) IPRange
}

type V4Searcher struct{}

func (s *V4Searcher) Search(ipStr string, table *IPTable) IPRange {
	ip := internal.ParseIPv4(ipStr)
	if ip == nil {
		return nil
	}
	idx := sort.Search(len(table.data), func(i int) bool {
		return table.data[i].GTE(ip)
	})

	if idx < len(table.data) && table.data[idx].Contains(ip) {
		return table.data[idx]
	}
	return nil
}

type V6Searcher struct{}

func (s *V6Searcher) Search(ipStr string, table *IPTable) IPRange {
	ip := internal.ParseIPv6(ipStr)
	idx := sort.Search(len(table.data), func(i int) bool {
		return table.data[i].GTE(ip)
	})

	if idx < len(table.data) && table.data[idx].Contains(ip) {
		return table.data[idx]
	}
	return nil
}
