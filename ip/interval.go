package ip

import (
	"encoding/binary"
	"net"
	"sort"
)

// Uint128 is a 16 bytes number, instead of big.Int which is much slower
type Uint128 struct {
	A uint64
	B uint64
}

func u128(bs []byte) *Uint128 {
	return &Uint128{A: binary.BigEndian.Uint64(bs[:8]), B: binary.BigEndian.Uint64(bs[8:])}
}

func (u *Uint128) Cmp(other *Uint128) int {
	if u.A > other.A {
		return 1
	} else if u.A < other.A {
		return -1
	} else if u.B > other.B {
		return 1
	} else if u.B < other.B {
		return -1
	}
	return 0
}

// Interval is a ip range
type Interval struct {
	StartStr   string
	EndStr     string
	CountryIdx int
	IspIdx     int
	ProvIdx    int
	CityIdx    int
	Number     int
	Contains   func(ip net.IP) bool
	Gte        func(ip net.IP) bool
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

func (lst *V4IntervalList) Len() int {
	return len(*lst)
}
func (lst *V4IntervalList) Swap(i, j int) {
	obj := *lst
	obj[i], obj[j] = obj[j], obj[i]
}
func (lst *V4IntervalList) Less(i, j int) bool {
	obj := *lst
	return obj[i].Cmp(obj[j]) < 0
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

func (lst *V6IntervalList) Len() int {
	return len(*lst)
}
func (lst *V6IntervalList) Swap(i, j int) {
	obj := *lst
	obj[i], obj[j] = obj[j], obj[i]
}
func (lst *V6IntervalList) Less(i, j int) bool {
	obj := *lst
	return obj[i].Cmp(obj[j]) < 0
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
