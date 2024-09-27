package ip

import (
	"bufio"
	"gommon/convert"
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
func ParseV4Range(line string, table *IPTable) *IPRange {
	return parseRange(line, table, true)
}

func ParseV6Range(line string, table *IPTable) *IPRange {
	return parseRange(line, table, false)
}

func parseRange(line string, table *IPTable, isV4 bool) *IPRange {
	if line[0] == '#' {
		return nil
	}
	vs := strings.Split(line, "|")
	if len(vs) != 7 {
		return nil
	}
	var low, high net.IP
	if isV4 {
		low, high = convert.IPStr2IPv4(vs[0]), convert.IPStr2IPv4(vs[1])
	} else {
		low, high = convert.IPStr2IPv6(vs[0]), convert.IPStr2IPv6(vs[1])
	}
	if low == nil || high == nil {
		return nil
	}
	if cmp(low, high) == 1 {
		low, high = high, low
	}
	return &IPRange{
		low:        low,
		high:       high,
		StartStr:   vs[0],
		EndStr:     vs[1],
		CountryIdx: table.countries.Append(vs[2]),
		IspIdx:     table.isps.Append(vs[3]),
		ProvIdx:    table.provs.Append(vs[4]),
		CityIdx:    table.cities.Append(vs[5]),
		NumberIdx:  table.numbers.Append(vs[6]),
	}
}
