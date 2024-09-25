package ip

import (
	"encoding/binary"
	"net"
)

func NewArray() *array {
	return &array{
		items:   make([]string, 0),
		hashmap: make(map[string]int),
	}
}

type array struct {
	items   []string
	hashmap map[string]int
}

func (a *array) append(ele string) int {
	idx, exists := a.hashmap[ele]
	if !exists {
		a.items = append(a.items, ele)
		idx = len(a.items) - 1
		a.hashmap[ele] = idx
	}
	return idx
}

func (a *array) get(idx int) string {
	return a.items[idx]
}

func uInt64Of(v6 string) uint64 {
	ip := net.ParseIP(v6)
	if ip == nil {
		return 0
	}
	v := ip.To16()
	if v == nil {
		return 0
	}
	return binary.BigEndian.Uint64(v)
}
