package verify

import (
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

func SpeciesCache() Cache {
	return Cache{get: species}
}
