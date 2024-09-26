package internal

import (
	"encoding/binary"
	"net"
)

type Int128 struct {
	H uint64
	L uint64
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
		H: binary.BigEndian.Uint64(v[0:8]),
		L: binary.BigEndian.Uint64(v[8:16]),
	}
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
