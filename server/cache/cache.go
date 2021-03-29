package cache

import (
	"reflect"
)

type Cache struct {
	memory *LRU
	disk   *BoltDb
}

func NewCache(capacity int, path string) *Cache {
	return &Cache{
		memory: NewLru(capacity),
		disk:   NewBoltDb(path),
	}
}

func (cache *Cache) Close() {
	cache.disk.Close()
}

func (cache *Cache) Get(key string, val interface{}) {
	cache.memory.Get(key, val)
	// TODO: There has to be a better way to do this
	if !reflect.ValueOf(val).Elem().IsNil() {
		return
	}

	cache.disk.Get(key, val)
}

func (cache *Cache) FastGet(key string, val interface{}) {
	cache.memory.Get(key, val)
}

func (cache *Cache) Put(key string, val interface{}) {
	cache.memory.Put(key, val)
	go cache.disk.Put(key, val)
}
