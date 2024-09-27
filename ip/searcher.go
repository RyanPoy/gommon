package ip

import (
	"gommon/convert"
	"net"
	"sort"
)

func SearchV4(ipStr string, ranges *IPRanges) *IPRange {
	ip := convert.IPStr2IPv4(ipStr)
	return search(ip, ranges)
}

func SearchV6(ipStr string, ranges *IPRanges) *IPRange {
	ip := convert.IPStr2IPv6(ipStr)
	return search(ip, ranges)
}

func search(ip net.IP, ranges *IPRanges) *IPRange {
	if ip == nil {
		return nil
	}
	idx := sort.Search(len(ranges.data), func(i int) bool {
		return cmp(ranges.data[i].low, ip) == 1 || cmp(ranges.data[i].high, ip) != -1
	})

	if idx < len(ranges.data) && ranges.data[idx].Contains(ip) {
		return ranges.data[idx]
	}
	return nil
}
