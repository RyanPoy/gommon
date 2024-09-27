package ip

import (
	"gommon/convert"
	"net"
	"sort"
)

func SearchV4(ipStr string, ranges IPRanges) *IPRange {
	ip := convert.IPStr2IPv4(ipStr)
	return search(ip, ranges)
}

func SearchV6(ipStr string, ranges IPRanges) *IPRange {
	ip := convert.IPStr2IPv6(ipStr)
	return search(ip, ranges)
}

func search(ip net.IP, ranges IPRanges) *IPRange {
	if ip == nil {
		return nil
	}
	idx := sort.Search(len(ranges), func(i int) bool {
		return cmp(ranges[i].low, ip) == 1 || cmp(ranges[i].high, ip) != -1
	})

	if idx < len(ranges) && ranges[idx].Contains(ip) {
		return ranges[idx]
	}
	return nil
}
