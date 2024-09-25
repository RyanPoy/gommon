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

func UInt128Of(v6 string) big.Int {
	ip := net.ParseIP(v6)
	if ip == nil {
		return *big.NewInt(0)
	}
	v := ip.To16()
	if v == nil {
		return *big.NewInt(0)
	}
	return *big.NewInt(0).SetBytes(v)
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
	cnt := strings.Count(v6, "::")
	if cnt == 0 || cnt > 1 {
		return v6
	}

	vs := strings.Split(v6, "::")
	if len(vs) != 2 {
		return v6
	}
	// 到这里一定是有1个::的字符串了
	cnt1, cnt2 := strings.Count(vs[0], ":"), strings.Count(vs[1], ":")
	if cnt1+cnt2+2 >= 8 {
		return v6
	}
	tmp := make([]string, 0)
	for i := 0; i < 8-cnt1-cnt2-2; i++ {
		tmp = append(tmp, "0000")
	}

	r1 := strings.Split(vs[0]+":"+strings.Join(tmp, ":")+":"+vs[1], ":")
	relt := make([]string, 0)
	for _, v := range r1 {
		l := len(v)
		if l == 0 {
			relt = append(relt, "0000")
		} else if l == 1 {
			relt = append(relt, "000"+v)
		} else if l == 2 {
			relt = append(relt, "00"+v)
		} else if l == 3 {
			relt = append(relt, "0"+v)
		} else {
			relt = append(relt, v)
		}
	}
	return strings.Join(relt, ":")
}
