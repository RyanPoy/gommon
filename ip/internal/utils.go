package internal

import (
	"encoding/binary"
	"net"
)

type Int128 struct {
	H uint64
	L uint64
}

func (i *Int128) Cmp(j *Int128) int {
	if i.H > j.H {
		return 1
	}
	if i.H < j.H {
		return -1
	}
	if i.L > j.L {
		return 1
	}
	if i.L < j.L {
		return -1
	}
	return 0
}

func NewI128From(ipStr string) *Int128 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	v := ip.To16()
	if v == nil {
		return nil
	}
	return &Int128{
		H: binary.BigEndian.Uint64(v[0:8]),
		L: binary.BigEndian.Uint64(v[8:16]),
	}
}

func ParseIPv6(ipStr string) net.IP {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	return ip.To16()
}

func ParseIPv4(ipStr string) net.IP {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	return ip.To4()
}
