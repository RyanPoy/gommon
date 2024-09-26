package extends

func NewArray() *Array {
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
