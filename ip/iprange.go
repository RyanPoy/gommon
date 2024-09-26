package ip

import (
	"gommon/ip/internal"
	"strings"
)

type IPRange interface {
	StartStr() string
	EndStr() string
	CountryIdx() int
	IspIdx() int
	ProvIdx() int
	CityIdx() int
	NumberIdx() int

	Cmp(IPRange) int
	GTE(ipStr string) bool
	Contains(ipValue interface{}) bool
}

type V4Range struct {
	low        uint32
	high       uint32
	startStr   string
	endStr     string
	countryIdx int
	ispIdx     int
	provIdx    int
	cityIdx    int
	numberIdx  int
}

func (r *V4Range) StartStr() string {
	return r.startStr
}
func (r *V4Range) EndStr() string {
	return r.endStr
}
func (r *V4Range) CountryIdx() int {
	return r.countryIdx
}
func (r *V4Range) IspIdx() int {
	return r.ispIdx
}
func (r *V4Range) ProvIdx() int {
	return r.provIdx
}
func (r *V4Range) CityIdx() int {
	return r.cityIdx
}
func (r *V4Range) NumberIdx() int {
	return r.numberIdx
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

func (r *V4Range) GTE(ipStr string) bool {
	ipv := uint32Of(ipStr)
	if r.low > ipv {
		return true
	}
	return r.high >= ipv
}

func (r *V4Range) Contains(ipValue interface{}) bool {
	ipv := *ipValue.(*uint32)
	return r.low <= ipv && ipv <= r.high
}

//
//
//

type V6Range struct {
	low        *internal.Int128
	high       *internal.Int128
	startStr   string
	endStr     string
	countryIdx int
	ispIdx     int
	provIdx    int
	cityIdx    int
	numberIdx  int
}

func (r *V6Range) StartStr() string {
	return r.startStr
}
func (r *V6Range) EndStr() string {
	return r.endStr
}
func (r *V6Range) CountryIdx() int {
	return r.countryIdx
}
func (r *V6Range) IspIdx() int {
	return r.ispIdx
}
func (r *V6Range) ProvIdx() int {
	return r.provIdx
}
func (r *V6Range) CityIdx() int {
	return r.cityIdx
}
func (r *V6Range) NumberIdx() int {
	return r.numberIdx
}
func (r *V6Range) Cmp(other IPRange) int {
	o := other.(*V6Range)

	cmpHigh := r.high.Cmp(o.high)
	if cmpHigh == 0 {
		return r.low.Cmp(o.low)
	}
	return cmpHigh
}

func (r *V6Range) GTE(ipStr string) bool {
	ipv := internal.FromIpv6(ipStr)

	if r.low.Cmp(ipv) == 1 {
		return true
	}
	return r.high.Cmp(ipv) != -1
}

func (r *V6Range) Contains(ipValue interface{}) bool {
	ipv := ipValue.(*internal.Int128)
	// 检查找到的index是否在原始区间内
	// 即：ipv <= r.low && ipv >= r.high
	return ipv.Cmp(r.low) != -1 && ipv.Cmp(r.high) != 1
}

func ParseV4Range(line string, table *IPTable) IPRange {
	if line[0] == '#' {
		return nil
	}
	vs := strings.Split(line, "|")
	if len(vs) != 7 {
		return nil
	}
	low := uint32Of(vs[0])
	high := uint32Of(vs[1])
	if low == 0 || high == 0 {
		return nil
	}
	if low > high {
		low, high = high, low
	}

	countryIdx := table.countries.Append(vs[2])
	ispIdx := table.isps.Append(vs[3])
	provIdx := table.provs.Append(vs[4])
	cityIdx := table.cities.Append(vs[5])
	numberIdx := table.numbers.Append(vs[6])

	return &V4Range{
		low:        low,
		high:       high,
		startStr:   vs[0],
		endStr:     vs[1],
		countryIdx: countryIdx,
		ispIdx:     ispIdx,
		provIdx:    provIdx,
		cityIdx:    cityIdx,
		numberIdx:  numberIdx,
	}
}

func ParseV6Range(line string, table *IPTable) IPRange {
	if line[0] == '#' {
		return nil
	}
	vs := strings.Split(line, "|")
	if len(vs) != 7 {
		return nil
	}

	low := internal.FromIpv6(vs[0])
	high := internal.FromIpv6(vs[1])
	if low == nil || high == nil {
		return nil
	}
	if low.Cmp(high) == 1 {
		low, high = high, low
	}

	countryIdx := table.countries.Append(vs[2])
	ispIdx := table.isps.Append(vs[3])
	provIdx := table.provs.Append(vs[4])
	cityIdx := table.cities.Append(vs[5])
	numberIdx := table.numbers.Append(vs[6])

	return &V6Range{
		low:        low,
		high:       high,
		startStr:   vs[0],
		endStr:     vs[1],
		countryIdx: countryIdx,
		ispIdx:     ispIdx,
		provIdx:    provIdx,
		cityIdx:    cityIdx,
		numberIdx:  numberIdx,
	}
}
