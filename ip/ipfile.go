package ip

import (
	"bufio"
	"gommon/ip/internal"
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
