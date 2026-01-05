package util

import (
	"sync"
)

func SlicesEqualUnordered[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	e := map[T]int{}
	for _, v := range a {
		e[v]++
	}

	for _, v := range b {
		if e[v] == 0 {
			return false
		}
		e[v]--
	}

	return true
}

func SlicesParMap[T, U any](s []T, f func(T) U) []U {
	res := make([]U, len(s))

	var wg sync.WaitGroup
	wg.Add(len(s))

	for i, x := range s {
		go func() {
			res[i] = f(x)
			wg.Done()
		}()
	}

	wg.Wait()

	return res
}
