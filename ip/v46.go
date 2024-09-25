package ip

import (
	"gommon/ip/internal"
	"sort"
	"strings"
)

type V46 struct {
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

type V46s struct {
	data      []V46
	countries *internal.Array
	isps      *internal.Array
	provs     *internal.Array
	cities    *internal.Array
	numbers   *internal.Array
}

func NewV46s(fpath string) (*V46s, error) {
	lines, err := internal.LoadFile(fpath)
	if err != nil {
		return nil, err
	}

	v46s := &V46s{
		data:      make([]V46, 0),
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
		var low, high uint64
		if strings.Contains(vs[0], ":") {
			vs[0] = internal.NormalizeV6(vs[0])
			vs[1] = internal.NormalizeV6(vs[1])
			low = internal.UInt64Of(vs[0])
			high = internal.UInt64Of(vs[1])
		} else {
			low = uint64(internal.UInt32Of(vs[0]))
			high = uint64(internal.UInt32Of(vs[1]))
		}
		if low == 0 || high == 0 {
			continue
		}
		if low > high {
			low, high = high, low
		}
		countryIdx := v46s.countries.Append(vs[2])
		ispIdx := v46s.isps.Append(vs[3])
		provIdx := v46s.provs.Append(vs[4])
		cityIdx := v46s.cities.Append(vs[5])
		numberIdx := v46s.numbers.Append(vs[6])

		v46s.data = append(v46s.data, V46{
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
	sort.Sort(v46s)
	return v46s, nil
}

func (v46s *V46s) Len() int {
	return len(v46s.data)
}

func (v46s *V46s) Swap(i, j int) {
	v46s.data[i], v46s.data[j] = v46s.data[j], v46s.data[i]
}

func (v46s *V46s) Less(i, j int) bool {
	o1, o2 := v46s.data[i], v46s.data[j]
	return o1.Low < o2.Low || (o1.Low == o2.Low && o1.High < o2.High)
}

func (v46s *V46s) Search(ipstr string) *V46 {
	var ipv uint64
	if strings.Contains(ipstr, ":") {
		ipstr = internal.NormalizeV6(ipstr)
		ipv = internal.UInt64Of(ipstr)
	} else {
		ipv = uint64(internal.UInt32Of(ipstr))
	}
	// 使用二分查找找到给定IP的合适位置
	idx := sort.Search(len(v46s.data), func(i int) bool {
		ip := v46s.data[i]
		if ip.Low > ipv {
			return true
		} else if ip.High < ipv {
			return false
		} else {
			return true
		}
	})

	// 检查找到的index是否在原始区间内
	if ipv >= v46s.data[idx].Low && ipv <= v46s.data[idx].High {
		return &v46s.data[idx]
	}
	return nil
}

func (v46s *V46s) StringOf(v4 *V46) string {
	return v4.StartStr + "|" +
		v4.EndStr + "|" +
		v46s.countries.Get(v4.CountryIdx) + "|" +
		v46s.isps.Get(v4.IspIdx) + "|" +
		v46s.provs.Get(v4.ProvIdx) + "|" +
		v46s.cities.Get(v4.CityIdx) + "|" +
		v46s.numbers.Get(v4.NumberIdx)
}
