package ip

import (
	"encoding/binary"
	"net"
	"sort"
)

// Interval is a ip range
type Interval interface {
	BaseInfo() *BaseInfo
	Cmp(o Interval) int
	Gte(ip net.IP) bool
	Contains(ip net.IP) bool
}

type BaseInfo struct {
	StartStr   string
	EndStr     string
	CountryIdx int
	IspIdx     int
	ProvIdx    int
	CityIdx    int
	Number     int
}

type V4Interval struct {
	baseInfo *BaseInfo
	Low      uint32
	High     uint32
}

func (i *V4Interval) BaseInfo() *BaseInfo {
	return i.baseInfo
}

func (i *V4Interval) Cmp(o Interval) int {
	other := o.(*V4Interval)
	if i.Low > other.High {
		return 1
	}
	if i.Low < other.High {
		return -1
	}
	return 0
}

func (i *V4Interval) Gte(ip net.IP) bool {
	ipv := binary.BigEndian.Uint32(ip)
	return i.Low > ipv || i.High >= ipv
}

func (i *V4Interval) Contains(ip net.IP) bool {
	ipv := binary.BigEndian.Uint32(ip)
	return i.Low <= ipv && ipv <= i.High
}

type V6Interval struct {
	baseInfo *BaseInfo
	Low      *Uint128
	High     *Uint128
}

func (i *V6Interval) BaseInfo() *BaseInfo {
	return i.baseInfo
}

func (i *V6Interval) Cmp(o Interval) int {
	other := o.(*V6Interval)
	cmpHigh := i.High.Cmp(other.High)
	if cmpHigh == 0 {
		return i.Low.Cmp(other.Low)
	}
	return cmpHigh
}

func (i *V6Interval) Contains(ip net.IP) bool {
	//return i.Low <= ipv && ipv <= i.High
	ipv := u128(ip)
	return i.Low.Cmp(ipv) <= 0 && ipv.Cmp(i.High) <= 0
}

func (i *V6Interval) Gte(ip net.IP) bool {
	ipv := u128(ip)
	return i.Low.Cmp(ipv) == 1 || i.High.Cmp(ipv) != -1
}

type IntervalList []Interval

func (lst IntervalList) Sort() {
	sort.Slice(lst, func(i, j int) bool {
		return lst[i].Cmp(lst[j]) < 0
	})
}

func (lst IntervalList) Search(ip net.IP) Interval {
	length := len(lst)
	idx := sort.Search(length, func(i int) bool {
		return lst[i].Gte(ip)
	})

	if idx < length && lst[idx].Contains(ip) {
		return lst[idx]
	}
	return nil
}

func (lst *IntervalList) Add(interval Interval) {
	*lst = append(*lst, interval)
}
