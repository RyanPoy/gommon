package ip

import (
	"gommon/ip/internal"
	"sort"
	"strings"
)

type V6 struct {
	Low      uint64
	High     uint64
	StartStr string
	EndStr   string

	CountryIdx int
	IspIdx     int
	ProvIdx    int
	CityIdx    int
	NumberIdx  int
}

type V6s struct {
	data      []V6
	countries *internal.Array
	isps      *internal.Array
	provs     *internal.Array
	cities    *internal.Array
	numbers   *internal.Array
}

func NewV6s(fpath string) (*V6s, error) {
	lines, err := internal.LoadFile(fpath)
	if err != nil {
		return nil, err
	}

	v6s := &V6s{
		data:      make([]V6, 0),
		countries: internal.NewArray(),
		isps:      internal.NewArray(),
		provs:     internal.NewArray(),
		cities:    internal.NewArray(),
		numbers:   internal.NewArray(),
	}

	// file format: start|end|country|isp|prov|city|line-number
	for _, line := range lines {
		if line[0] == '#' {
			continue
		}
		vs := strings.Split(line, "|")
		if len(vs) != 7 {
			continue
		}
		vs[0] = internal.NormalizeV6(vs[0])
		vs[1] = internal.NormalizeV6(vs[1])

		low := internal.UInt64Of(vs[0])
		high := internal.UInt64Of(vs[1])
		if low == 0 || high == 0 {
			continue
		}
		if low > high {
			low, high = high, low
		}
		countryIdx := v6s.countries.Append(vs[2])
		ispIdx := v6s.isps.Append(vs[3])
		provIdx := v6s.provs.Append(vs[4])
		cityIdx := v6s.cities.Append(vs[5])
		numberIdx := v6s.numbers.Append(vs[6])

		v6s.data = append(v6s.data, V6{
			Low:        low,
			High:       high,
			StartStr:   vs[0],
			EndStr:     vs[1],
			CountryIdx: countryIdx,
			IspIdx:     ispIdx,
			ProvIdx:    provIdx,
			CityIdx:    cityIdx,
			NumberIdx:  numberIdx,
		})
	}
	sort.Sort(v6s)
	return v6s, nil
}

func (v6s *V6s) Len() int {
	return len(v6s.data)
}

func (v6s *V6s) Swap(i, j int) {
	v6s.data[i], v6s.data[j] = v6s.data[j], v6s.data[i]
}

func (v6s *V6s) Less(i, j int) bool {
	o1, o2 := v6s.data[i], v6s.data[j]
	return o1.Low < o2.Low || (o1.Low == o2.Low && o1.High < o2.High)
}

func (v6s *V6s) Search(ipstr string) *V6 {
	ipstr = internal.NormalizeV6(ipstr)
	ipv := internal.UInt64Of(ipstr)

	// 使用二分查找找到给定IP的合适位置
	idx := sort.Search(len(v6s.data), func(i int) bool {
		ip := v6s.data[i]
		if ip.Low > ipv {
			return true
		} else if ip.High < ipv {
			return false
		} else {
			return true
		}
	})

	// 检查找到的index是否在原始区间内
	if ipv >= v6s.data[idx].Low && ipv <= v6s.data[idx].High {
		return &v6s.data[idx]
	}
	return nil
}

func (v6s *V6s) StringOf(v6 *V6) string {
	return v6.StartStr + "|" +
		v6.EndStr + "|" +
		v6s.countries.Get(v6.CountryIdx) + "|" +
		v6s.isps.Get(v6.IspIdx) + "|" +
		v6s.provs.Get(v6.ProvIdx) + "|" +
		v6s.cities.Get(v6.CityIdx) + "|" +
		v6s.numbers.Get(v6.NumberIdx)
}
