package ip

import (
	"bufio"
	"encoding/binary"
	"gommon/ip/internal"
	"net"
	"os"
	"strings"
)

func LoadFile(fpath string) ([]string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]string, 0)
	for scanner := bufio.NewScanner(f); scanner.Scan(); {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
func ParseV4Range(line string, table *IPTable) IPRange {
	if line[0] == '#' {
		return nil
	}
	vs := strings.Split(line, "|")
	if len(vs) != 7 {
		return nil
	}
	low := u32(vs[0])
	high := u32(vs[1])
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

	low := bigInt(vs[0])
	high := bigInt(vs[1])
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

func bigInt(ipStr string) *internal.Int128 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	v := ip.To16()
	if v == nil {
		return nil
	}
	return &internal.Int128{
		H: binary.BigEndian.Uint64(v[0:8]),
		L: binary.BigEndian.Uint64(v[8:16]),
	}
}
func u32(v4 string) uint32 {
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
