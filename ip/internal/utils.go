package internal

import (
	"bufio"
	"encoding/binary"
	"math/big"
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

func UInt128Of(v6 string) *big.Int {
	ip := net.ParseIP(v6)
	if ip == nil {
		return nil
	}
	v := ip.To16()
	if v == nil {
		return nil
	}
	return big.NewInt(0).SetBytes(v)
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
	result := make([]string, 0, 8)
	for _, segment := range leftPart {
		result = append(result, strings.Repeat("0", 4-len(segment))+segment)
	}

	for i := 0; i < zerosToAdd; i++ {
		result = append(result, "0000")
	}

	for _, segment := range rightPart {
		result = append(result, strings.Repeat("0", 4-len(segment))+segment)
	}

	return strings.Join(result, ":")
}
