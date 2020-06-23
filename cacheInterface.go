package caches

const (
	cacheSize = 64
	addressMaxNumber = 4294967295 //2^32-1
)

//Cache - an interface for cache
type Cache interface {
	//Fetch - getting data and if was a hit, by passing address
	Fetch(address uint32) (byte, bool)

	//Store - updating data and if was a hit, by passing address
	Store(address uint32, newData byte) bool
}