package extends

type Int128 struct {
	H uint64
	L uint64
}

func (i *Int128) Cmp(j *Int128) int {
	if i.H > j.H {
		return 1
	}
	if i.H < j.H {
		return -1
	}
	if i.L > j.L {
		return 1
	}
	if i.L < j.L {
		return -1
	}
	return 0
}
