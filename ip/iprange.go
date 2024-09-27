package ip

import (
	"bytes"
	"encoding/binary"
	"net"
)

//type UInt128 []byte

var cmp = bytes.Compare

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
	low        []byte
	high       []byte
}

func (r *V6Range) OriginData() *OriginData {
	return r.originData
}

func (r *V6Range) Cmp(other IPRange) int {
	o := other.(*V6Range)

	cmpHigh := cmp(r.high, o.high)
	if cmpHigh == 0 {
		return cmp(r.low, o.low)
	}
	return cmpHigh
}

func (r *V6Range) GTE(ip net.IP) bool {
	if cmp(r.low, ip) == 1 {
		return true
	}
	return cmp(r.high, ip) != -1
}

func (r *V6Range) Contains(ip net.IP) bool {
	// 检查找到的index是否在原始区间内
	// 即：ipv <= r.low && ipv >= r.high
	return cmp(ip, r.low) != -1 && cmp(ip, r.high) != 1
}

func NewV4Range(
	low, high net.IP,
	startStr, endStr string,
	countryIdx, ispIdx, provIdx, cityIdx, numberIdx int) *V4Range {
	return &V4Range{
		low:  binary.BigEndian.Uint32(low),
		high: binary.BigEndian.Uint32(high),
		originData: &OriginData{
			StartStr:   startStr,
			EndStr:     endStr,
			CountryIdx: countryIdx,
			IspIdx:     ispIdx,
			ProvIdx:    provIdx,
			CityIdx:    cityIdx,
			NumberIdx:  numberIdx,
		},
	}
}

func NewV6Range(low, high net.IP,
	startStr, endStr string,
	countryIdx, ispIdx, provIdx, cityIdx, numberIdx int) *V6Range {
	return &V6Range{
		low:  low,
		high: high,
		originData: &OriginData{
			StartStr:   startStr,
			EndStr:     endStr,
			CountryIdx: countryIdx,
			IspIdx:     ispIdx,
			ProvIdx:    provIdx,
			CityIdx:    cityIdx,
			NumberIdx:  numberIdx,
		},
	}
}
