package Caches

const (
	cacheSize = 64
	addressMaxNumber = 4294967295 //2^32-1
)

type Cache interface {
	//getting data and if was a hit, by passing address
	Fetch(address uint32) (byte, bool)
	//updating data and if was a hit, by passing address
	Store(address uint32, newData byte) bool
}