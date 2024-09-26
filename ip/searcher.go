package ip

import (
	"net"
	"sort"
)

type Searcher interface {
	Search(ipStr string, table *IPTable) IPRange
}

type V4Searcher struct{}

func (s *V4Searcher) Search(ipStr string, table *IPTable) IPRange {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	ipv := ip.To4()
	if ipv == nil {
		return nil
	}
	idx := sort.Search(len(table.data), func(i int) bool {
		return table.data[i].GTE(ipv)
	})

	if idx < len(table.data) && table.data[idx].Contains(ipv) {
		return table.data[idx]
	}
	return nil
}

type V6Searcher struct{}

func (s *V6Searcher) Search(ipStr string, table *IPTable) IPRange {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	ipv := ip.To16()
	if ipv == nil {
		return nil
	}
	idx := sort.Search(len(table.data), func(i int) bool {
		return table.data[i].GTE(ipv)
	})

	if idx < len(table.data) && table.data[idx].Contains(ipv) {
		return table.data[idx]
	}
	return nil
}
