package cache

type Cache map[string]*any

var cache = make(Cache)

func LookUp(key string) any {
	return cache[key]
}

func Put(key string, value any) {
	if v := cache[key]; v != nil {
		if v == value {
			return
		}
		delete(cache, key)
		cache[key] = &value
		return
	}
	cache[key] = &value
}
