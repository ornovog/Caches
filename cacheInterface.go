package Caches

const cacheSize = 64

type Cache interface {
	//getting data and if was a hit, by passing address
	GetData(address uint32) (float64, bool)
}