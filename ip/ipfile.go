package ip

import (
	"bufio"
	"gommon/convert"
	"gommon/extends"
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
	return parseRange(line, table, true)
}

func ParseV6Range(line string, table *IPTable) IPRange {
	return parseRange(line, table, false)
}

func parseRange(line string, table *IPTable, isV4 bool) IPRange {
	if line[0] == '#' {
		return nil
	}
	vs := strings.Split(line, "|")
	if len(vs) != 7 {
		return nil
	}
	low, high := parseLowAndHigh(vs[0], vs[1], isV4)
	if low == nil || high == nil {
		return nil
	}

	if isV4 {
		return &V4Range{
			low:        low.(uint32),
			high:       high.(uint32),
			startStr:   vs[0],
			endStr:     vs[1],
			countryIdx: table.countries.Append(vs[2]),
			ispIdx:     table.isps.Append(vs[3]),
			provIdx:    table.provs.Append(vs[4]),
			cityIdx:    table.cities.Append(vs[5]),
			numberIdx:  table.numbers.Append(vs[6]),
		}
	}
	return &V6Range{
		low:        low.(*extends.Int128),
		high:       high.(*extends.Int128),
		startStr:   vs[0],
		endStr:     vs[1],
		countryIdx: table.countries.Append(vs[2]),
		ispIdx:     table.isps.Append(vs[3]),
		provIdx:    table.provs.Append(vs[4]),
		cityIdx:    table.cities.Append(vs[5]),
		numberIdx:  table.numbers.Append(vs[6]),
	}
}

func parseLowAndHigh(lowStr, highStr string, isV4 bool) (interface{}, interface{}) {
	if isV4 {
		low := convert.IPStr2Uint32(lowStr)
		high := convert.IPStr2Uint32(highStr)
		if low == 0 || high == 0 {
			return nil, nil
		}
		if low > high {
			low, high = high, low
		}
		return low, high
	} else {
		low := convert.IPStr2Int128(lowStr)
		high := convert.IPStr2Int128(highStr)
		if low == nil || high == nil {
			return nil, nil
		}
		if low.Cmp(high) == 1 {
			low, high = high, low
		}
		return low, high
	}
}
