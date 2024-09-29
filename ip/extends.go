package ip

import "encoding/binary"

func array() *Array {
	return &Array{
		items:   make([]string, 0),
		hashmap: make(map[string]int),
	}
}

type Array struct {
	items   []string
	hashmap map[string]int
}

func (a *Array) Append(ele string) int {
	idx, exists := a.hashmap[ele]
	if !exists {
		a.items = append(a.items, ele)
		idx = len(a.items) - 1
		a.hashmap[ele] = idx
	}
	return idx
}

func (a *Array) Get(idx int) string {
	return a.items[idx]
}

// Uint128 is a 16 bytes number, instead of big.Int which is much slower
type Uint128 struct {
	A uint64
	B uint64
}

func u128(bs []byte) *Uint128 {
	return &Uint128{A: binary.BigEndian.Uint64(bs[:8]), B: binary.BigEndian.Uint64(bs[8:])}
}

func (u *Uint128) Cmp(other *Uint128) int {
	if u.A > other.A {
		return 1
	} else if u.A < other.A {
		return -1
	} else if u.B > other.B {
		return 1
	} else if u.B < other.B {
		return -1
	}
	return 0
}
