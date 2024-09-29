package ip

import (
	"bytes"
	"encoding/binary"
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
	Cmp        func(other *IPRange) int
	Contains   func(ip net.IP) bool
	Gte        func(ip net.IP) bool
}

type V4Range struct {
	IPRange
	Low  *uint32
	High *uint32
}

func (r *V4Range) Cmp(other *V4Range) int {
	if *r.Low > *other.High {
		return 1
	}
	if *r.Low < *other.High {
		return -1
	}
	return 0
	//low, high := binary.BigEndian.Uint32(r.Low), binary.BigEndian.Uint32(other.High)
	//if low > high {
	//	return 1
	//}
	//if low < high {
	//	return -1
	//}
	//return 0
}
func (r *V4Range) Contains(ip net.IP) bool {
	ipv := binary.BigEndian.Uint32(ip)
	return *r.Low <= ipv && ipv <= *r.High

	//ipv := binary.BigEndian.Uint32(ip)
	//low, high := binary
	//return *r.Low <= ipv && ipv <= *r.High
}

func (r *V4Range) Gte(ip net.IP) bool {
	//cmp(r.Low, ip) == 1 || cmp(r.High, ip) != -1
	ipv := binary.BigEndian.Uint32(ip)
	return *r.Low > ipv || *r.High >= ipv
}

type V6Range struct {
	IPRange
	Low  net.IP
	High net.IP
}

func (r *V6Range) Cmp(other *V6Range) int {
	cmpHigh := cmp(r.High, other.High)
	if cmpHigh == 0 {
		return cmp(r.Low, other.Low)
	}
	return cmpHigh
}

func (r *V6Range) Contains(ip net.IP) bool {
	//return r.Low <= ipv && ipv <= r.High
	return cmp(r.Low, ip) <= 0 && cmp(ip, r.High) <= 0
}

func (r *V6Range) Gte(ip net.IP) bool {
	return cmp(r.Low, ip) == 1 || cmp(r.High, ip) != -1
}

// /
// /
// /
type V4Ranges []*V4Range

func (rs *V4Ranges) Len() int {
	return len(*rs)
}
func (rs *V4Ranges) Swap(i, j int) {
	obj := *rs
	obj[i], obj[j] = obj[j], obj[i]
}
func (rs *V4Ranges) Less(i, j int) bool {
	obj := *rs
	return obj[i].Cmp(obj[j]) < 0
}

type V6Ranges []*V6Range

func (rs *V6Ranges) Len() int {
	return len(*rs)
}
func (rs *V6Ranges) Swap(i, j int) {
	obj := *rs
	obj[i], obj[j] = obj[j], obj[i]
}
func (rs *V6Ranges) Less(i, j int) bool {
	obj := *rs
	return obj[i].Cmp(obj[j]) < 0
}
