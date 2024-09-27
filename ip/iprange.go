package ip

import (
	"encoding/binary"
	"gommon/extends"
	"net"
)

type OriginData struct {
	StartStr   string
	EndStr     string
	CountryIdx int
	IspIdx     int
	ProvIdx    int
	CityIdx    int
	NumberIdx  int
}

type IPRange interface {
	OriginData() *OriginData
	Cmp(IPRange) int
	GTE(ip net.IP) bool
	Contains(ip net.IP) bool
}

type V4Range struct {
	originData *OriginData
	low        uint32
	high       uint32
}

func (r *V4Range) OriginData() *OriginData {
	return r.originData
}

func (r *V4Range) Cmp(other IPRange) int {
	o := other.(*V4Range)
	if r.low > o.low {
		return 1
	} else if r.low < o.low {
		return -1
	} else if r.high > o.high {
		return 1
	} else if r.high < o.high {
		return -1
	}
	return 0
}

func (r *V4Range) GTE(ip net.IP) bool {
	ipv := binary.BigEndian.Uint32(ip)
	if r.low > ipv {
		return true
	}
	return r.high >= ipv
}

func (r *V4Range) Contains(ip net.IP) bool {
	ipv := binary.BigEndian.Uint32(ip)
	return r.low <= ipv && ipv <= r.high
}

//
//
//

type V6Range struct {
	originData *OriginData
	low        *extends.Int128
	high       *extends.Int128
}

func (r *V6Range) OriginData() *OriginData {
	return r.originData
}

func (r *V6Range) Cmp(other IPRange) int {
	o := other.(*V6Range)

	cmpHigh := r.high.Cmp(o.high)
	if cmpHigh == 0 {
		return r.low.Cmp(o.low)
	}
	return cmpHigh
}

func (r *V6Range) GTE(ip net.IP) bool {
	ipv := &extends.Int128{
		H: binary.BigEndian.Uint64(ip[0:8]),
		L: binary.BigEndian.Uint64(ip[8:16]),
	}
	if r.low.Cmp(ipv) == 1 {
		return true
	}
	return r.high.Cmp(ipv) != -1
}

func (r *V6Range) Contains(ip net.IP) bool {
	ipv := &extends.Int128{
		H: binary.BigEndian.Uint64(ip[0:8]),
		L: binary.BigEndian.Uint64(ip[8:16]),
	}
	// 检查找到的index是否在原始区间内
	// 即：ipv <= r.low && ipv >= r.high
	return ipv.Cmp(r.low) != -1 && ipv.Cmp(r.high) != 1
}
