package convert

import (
	"net"
)

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
