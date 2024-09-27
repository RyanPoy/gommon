package convert

import (
	"encoding/binary"
	"net"
)

var IPStr2Int128 = IPStr2IPv6

func IPStr2Uint32(ipStr string) *uint32 {
	v := IPStr2IPv4(ipStr)
	if v == nil {
		return nil
	}
	r := binary.BigEndian.Uint32(v)
	return &r
}
func IPStr2IPv6(ipStr string) net.IP {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	return ip.To16()
}

func IPStr2IPv4(ipStr string) net.IP {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	return ip.To4()
}
