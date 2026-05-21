package util

import (
	"sync"
)

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
