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

func DOICache(noVerifyReq bool) Cache {
	return Cache{get: func(s string, ctx context.Context) error {
		return doi(s, ctx, noVerifyReq)
	}}
}

func ORCIDCache(noVerifyReq bool) Cache {
	return Cache{get: func(s string, ctx context.Context) error {
		return orcid(s, ctx, noVerifyReq)
	}}
}

func SpeciesCache(noVerifyReq bool) Cache {
	return Cache{get: func(s string, ctx context.Context) error {
		return species(s, ctx, noVerifyReq)
	}}
}

func BioPortalCache(o types.Onto, bioPortalAPIKey string, noVerifyReq bool) Cache {
	return Cache{get: func(s string, ctx context.Context) error {
		return bioPortal(o, s, ctx, bioPortalAPIKey, noVerifyReq)
	}}
}

func OLS4Cache(o types.Onto, noVerifyReq bool) Cache {
	return Cache{get: func(s string, ctx context.Context) error {
		return ols4(o, s, ctx, noVerifyReq)
	}}
}
