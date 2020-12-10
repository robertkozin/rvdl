package cache

import (
	"container/list"
	"reflect"
)

// TODO: I don't like how the code in this file looks

type LRU struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
}

type entry struct {
	key string
	val interface{}
}

func NewLru(size int) *LRU {
	return &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element, size),
	}
}

func (lru *LRU) Get(key string, val interface{}) {
	if el, ok := lru.items[key]; ok {
		lru.evictList.MoveToFront(el)

		ret := el.Value.(*entry).val
		reflect.ValueOf(val).Elem().Set(reflect.ValueOf(ret))
	}
}

func (lru *LRU) Put(key string, val interface{}) {
	if el, ok := lru.items[key]; ok {
		el.Value.(*entry).val = val
		lru.evictList.MoveToFront(el)
		return
	}

	if len(lru.items) >= lru.size {
		delete(lru.items, lru.evictList.Remove(lru.evictList.Back()).(*entry).key)
	}

	lru.items[key] = lru.evictList.PushFront(&entry{key, val})
}
