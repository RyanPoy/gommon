package internal

type Int128 struct {
	h uint64
	l uint64
}

func (i *Int128) Cmp(j *Int128) int {
	if i.h > j.h {
		return 1
	} else if i.h < j.h {
		return -1
	} else if i.l > j.l {
		return 1
	} else if i.l < j.l {
		return -1
	}
	return 0
}
