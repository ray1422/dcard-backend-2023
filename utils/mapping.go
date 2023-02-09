package utils

// Map is just like python's map which takes a function and a slice.
func Map[t any, s any](fun func(t) s, inp []t) (ret []s) {
	ret = make([]s, len(inp))
	for i, v := range inp {
		ret[i] = fun(v)
	}
	return
}

// MapInPlace is like map but modify the object directly without copy.
func MapInPlace[t any](fun func(*t), inp []t) {
	for i := range inp {
		fun(&inp[i])
	}
}
