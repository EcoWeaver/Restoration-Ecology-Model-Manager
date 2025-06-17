package verify

import (
	"sync"
)

type Cache struct {
	get   func(string) error
	cache sync.Map
}

func (self *Cache) Verify(key string) error {
	if _, ok := self.cache.Load(key); ok {
		return nil
	}

	err := self.get(key)
	if err == nil {
		self.cache.Store(key, struct{}{})
	}
	return err
}

func DoiCache() Cache {
	return Cache{get: doi}
}

func OrcidCache() Cache {
	return Cache{get: orcid}
}

func SpeciesCache() Cache {
	return Cache{get: species}
}
