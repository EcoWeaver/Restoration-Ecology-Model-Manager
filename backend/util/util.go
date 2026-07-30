package util

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

func GetAndUnmarshal[T any](url string, ctx context.Context) (T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return *new(T), err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return *new(T), err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return *new(T), err
	}

	var r T

	if err = json.Unmarshal(body, &r); err != nil {
		return *new(T), err
	}

	return r, nil
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
