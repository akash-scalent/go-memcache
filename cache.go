package main

import "github.com/bradfitz/gomemcache/memcache"

func newCache(address string) *memcache.Client {
	cache := memcache.New(address)
	return cache
}
