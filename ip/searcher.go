package ip

import (
	"gommon/convert"
	"net"
	"sort"
)

func SearchV4(ipStr string, table *IPTable) *IPRange {
	ip := convert.IPStr2IPv4(ipStr)
	return search(ip, table)
}

func SearchV6(ipStr string, table *IPTable) *IPRange {
	ip := convert.IPStr2IPv6(ipStr)
	return search(ip, table)
}

func search(ip net.IP, table *IPTable) *IPRange {
	if ip == nil {
		return nil
	}
	idx := sort.Search(len(table.data), func(i int) bool {
		return cmp(table.data[i].low, ip) == 1 || cmp(table.data[i].high, ip) != -1
	})

	if idx < len(table.data) && table.data[idx].Contains(ip) {
		return table.data[idx]
	}
	return nil
}
