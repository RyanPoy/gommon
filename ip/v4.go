package ip

import (
	"encoding/binary"
	"net"
	"sort"
	"strings"
)

type V4 struct {
	Low      uint32
	High     uint32
	StartStr string
	EndStr   string

	CountryIdx int
	IspIdx     int
	ProvIdx    int
	CityIdx    int
	NumberIdx  int
}

type V4s struct {
	data      []V4
	countries *array
	isps      *array
	provs     *array
	cities    *array
	numbers   *array
}

func NewV4s(fpath string) (*V4s, error) {
	lines, err := loadFile(fpath)
	if err != nil {
		return nil, err
	}

	v4s := &V4s{
		data:      make([]V4, 0),
		countries: NewArray(),
		isps:      NewArray(),
		provs:     NewArray(),
		cities:    NewArray(),
		numbers:   NewArray(),
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
		low := uInt32Of(vs[0])
		high := uInt32Of(vs[1])
		if low == 0 || high == 0 {
			continue
		}
		if low > high {
			low, high = high, low
		}
		countryIdx := v4s.countries.append(vs[2])
		ispIdx := v4s.isps.append(vs[3])
		provIdx := v4s.provs.append(vs[4])
		cityIdx := v4s.cities.append(vs[5])
		numberIdx := v4s.numbers.append(vs[6])

		v4s.data = append(v4s.data, V4{
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
	sort.Sort(v4s)
	return v4s, nil
}

func uInt32Of(v4 string) uint32 {
	ip := net.ParseIP(v4)
	if ip == nil {
		return 0
	}
	v := ip.To4()
	if v == nil {
		return 0
	}
	return binary.BigEndian.Uint32(v)
}

func (v4s *V4s) Len() int {
	return len(v4s.data)
}

func (v4s *V4s) Swap(i, j int) {
	v4s.data[i], v4s.data[j] = v4s.data[j], v4s.data[i]
}

func (v4s *V4s) Less(i, j int) bool {
	return less(v4s.data[i], v4s.data[j])
}

func less(o1, o2 V4) bool {
	return o1.Low < o2.Low || (o1.Low == o2.Low && o1.High < o2.High)
}
func (v4s *V4s) Search(ipstr string) *V4 {
	ipv := uInt32Of(ipstr)

	// 使用二分查找找到给定IP的合适位置
	idx := sort.Search(len(v4s.data), func(i int) bool {
		ip := v4s.data[i]
		if ip.Low > ipv {
			return true
		} else if ip.High < ipv {
			return false
		} else {
			return true
		}
	})

	// 检查找到的index是否在原始区间内
	if ipv >= v4s.data[idx].Low && ipv <= v4s.data[idx].High {
		return &v4s.data[idx]
	}
	return nil
}

func (v4s *V4s) StringOf(v4 *V4) string {
	return v4.StartStr + "|" +
		v4.EndStr + "|" +
		v4s.countries.get(v4.CountryIdx) + "|" +
		v4s.isps.get(v4.IspIdx) + "|" +
		v4s.provs.get(v4.ProvIdx) + "|" +
		v4s.cities.get(v4.CityIdx) + "|" +
		v4s.numbers.get(v4.NumberIdx)
}
