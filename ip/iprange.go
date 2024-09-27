package ip

import (
	"bytes"
	"net"
)

var cmp = bytes.Compare

type IPRange struct {
	StartStr   string
	EndStr     string
	CountryIdx int
	IspIdx     int
	ProvIdx    int
	CityIdx    int
	NumberIdx  int
	low        net.IP
	high       net.IP
}

func (r *IPRange) Cmp(other *IPRange) int {
	cmpHigh := cmp(r.high, other.high)
	if cmpHigh == 0 {
		return cmp(r.low, other.low)
	}
	return cmpHigh
}

func (r *IPRange) Contains(ip net.IP) bool {
	//return r.low <= ipv && ipv <= r.high
	return cmp(r.low, ip) <= 0 && cmp(ip, r.high) <= 0
}
