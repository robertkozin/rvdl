package cache

type Kv interface {
	Get(key string, val interface{})
	Put(key string, val interface{})
}
