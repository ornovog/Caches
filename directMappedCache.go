package Caches

const (
	indexBits = 63
	tagBits = 4294967232
)

type DMCacheLine struct {
	tag uint32
	data float64
}

type directMappedCache struct{
	storage [cacheSize]DMCacheLine
	mM *mainMemory
}

func (dMC *directMappedCache) GetData(address uint32) (float64, bool){
	index := address & indexBits
	tag := address & tagBits
	line := dMC.storage[index]

	if line.tag == tag{
		return line.data, true
	}else {
		data := dMC.mM.Fetch(address)
		dMC.storage[index].data = data
		dMC.storage[index].tag = tag

		return data, false
	}
}