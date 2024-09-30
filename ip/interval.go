package ip

import (
	"encoding/binary"
	"net"
	"sort"
)

// Interval is a ip range
type Interval struct {
	StartStr   string
	EndStr     string
	CountryIdx int
	IspIdx     int
	ProvIdx    int
	CityIdx    int
	Number     int

	//Cmp      func(other *Interval) int
	Contains func(ip net.IP) bool
}

type V4Interval struct {
	Interval
	Low  uint32
	High uint32
}

func (i *V4Interval) Cmp(other *V4Interval) int {
	if i.Low > other.High {
		return 1
	}
	if i.Low < other.High {
		return -1
	}
	return 0
}
func (i *V4Interval) Contains(ip net.IP) bool {
	ipv := binary.BigEndian.Uint32(ip)
	return i.Low <= ipv && ipv <= i.High
}

type V6Interval struct {
	Interval
	Low  *Uint128
	High *Uint128
}

func (i *V6Interval) Cmp(other *V6Interval) int {
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

type V4IntervalList []*V4Interval

func (lst V4IntervalList) Sort() {
	obj := lst
	sort.Slice(lst, func(i, j int) bool {
		return obj[i].Cmp(obj[j]) < 0
	})
}
func (lst V4IntervalList) Search(ip net.IP) *V4Interval {
	length := len(lst)
	idx := sort.Search(length, func(i int) bool {
		ipv := binary.BigEndian.Uint32(ip)
		return lst[i].Low > ipv || lst[i].High >= ipv
	})

	if idx < length && lst[idx].Contains(ip) {
		return lst[idx]
	}
	return nil
}

type V6IntervalList []*V6Interval

func (lst V6IntervalList) Sort() {
	obj := lst
	sort.Slice(lst, func(i, j int) bool {
		return obj[i].Cmp(obj[j]) < 0
	})
}

func (lst V6IntervalList) Search(ip net.IP) *V6Interval {
	length := len(lst)
	idx := sort.Search(length, func(i int) bool {
		ipv := u128(ip)
		return lst[i].Low.Cmp(ipv) == 1 || lst[i].High.Cmp(ipv) != -1
	})

	if idx < length && lst[idx].Contains(ip) {
		return lst[idx]
	}
	return nil
}
