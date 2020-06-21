package Caches

const (
	cacheSize = 64
	addressMaxNumber = 4294967295 //2^32-1
)

type Cache interface {
	//getting data and if was a hit, by passing address
	GetData(address uint32) (float64, bool)
}