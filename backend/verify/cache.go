package verify

import (
	"backend/types"
	"context"
	"sync"
)

type Cache struct {
	get   func(string, context.Context) error
	cache sync.Map
}

func (self *Cache) Verify(s string, ctx context.Context) error {
	if _, ok := self.cache.Load(s); ok {
		return nil
	}

	err := self.get(s, ctx)
	if err == nil {
		self.cache.Store(s, struct{}{})
	}
	return err
}

func DOICache() Cache {
	return Cache{get: doi}
}

func ORCIDCache() Cache {
	return Cache{get: orcid}
}

func SpeciesCache() Cache {
	return Cache{get: species}
}

func OLS4Cache(o types.Onto) Cache {
	return Cache{get: func(s string, ctx context.Context) error {
		return ols4(o, s, ctx)
	}}
}
