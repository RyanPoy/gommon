package internal

import (
	"bufio"
	"encoding/binary"
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

func UInt128Of(v6 string) *Int128 {
	ip := net.ParseIP(v6)
	if ip == nil {
		return nil
	}
	v := ip.To16()
	if v == nil {
		return nil
	}
	return &Int128{
		h: binary.BigEndian.Uint64(v[0:8]),
		l: binary.BigEndian.Uint64(v[8:16]),
	}
}

func UInt32Of(v4 string) uint32 {
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

func NormalizeV6(v6 string) string {
	if v6 == "::" {
		return "0000:0000:0000:0000:0000:0000:0000:0000"
	}
	if len(v6) == 39 { // 认为这个长度就是正确的ipv6长度
		return v6
	}

	cnt := strings.Count(v6, "::")
	if cnt != 1 {
		return v6
	}
	// 到这里一定是有1个::的字符串了

	// "::" 分割左右部分
	parts := strings.Split(v6, "::")
	leftPart := strings.Split(parts[0], ":")
	rightPart := strings.Split(parts[1], ":")

	// 计算需要补全的 "0000" 的数量
	zerosToAdd := 8 - len(leftPart) - len(rightPart)

	// 构建补齐后的结果
	result := make([]string, 8)
	idx := 0
	for i := range leftPart {
		result[idx] = zeroPad(leftPart[i])
		idx++
	}

	for i := 0; i < zerosToAdd; i++ {
		result[idx] = "0000"
		idx++
	}

	for i := range rightPart {
		result[idx] = zeroPad(rightPart[i])
		idx++
	}

	return strings.Join(result, ":")
}

func zeroPad(s string) string {
	switch len(s) {
	case 0:
		return "0000"
	case 1:
		return "000" + s
	case 2:
		return "00" + s
	case 3:
		return "0" + s
	default:
		return s
	}
}
