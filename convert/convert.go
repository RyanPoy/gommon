package convert

import (
	"encoding/binary"
	"gommon/extends"
	"net"
)

func IPStr2Int128(ipStr string) *extends.Int128 {
	v := IPStr2IPv6(ipStr)
	if v == nil {
		return nil
	}
	return &extends.Int128{
		H: binary.BigEndian.Uint64(v[0:8]),
		L: binary.BigEndian.Uint64(v[8:16]),
	}
}

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
