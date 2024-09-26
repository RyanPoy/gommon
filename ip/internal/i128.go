package internal

import (
	"encoding/binary"
	"net"
)

type Int128 struct {
	h uint64
	l uint64
}

func FromIpv6(v6 string) *Int128 {
	ip := net.ParseIP(v6)
	if ip == nil {
		return nil
	}
	v := ip.To16()
	if v == nil {
		return nil
	}
	return &Int128{
		h: binary.BigEndian.Uint64(v[0:8]),
		l: binary.BigEndian.Uint64(v[8:16]),
	}
}

func NewInt128(v6 string) *Int128 {
	ip := net.ParseIP(v6)
	if ip == nil {
		return nil
	}
	v := ip.To16()
	if v == nil {
		return nil
	}
	return &Int128{
		h: binary.BigEndian.Uint64(v[0:8]),
		l: binary.BigEndian.Uint64(v[8:16]),
	}
}

func (i *Int128) Cmp(j *Int128) int {
	if i.h > j.h {
		return 1
	}
	if i.h < j.h {
		return -1
	}
	if i.l > j.l {
		return 1
	}
	if i.l < j.l {
		return -1
	}
	return 0
}
