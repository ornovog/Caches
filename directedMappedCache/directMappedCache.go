package directedMappedCache

import (
	"Caches"
	"Caches/mainMemory"
)

const (
	indexBits = caches.CacheSize - 1
	tagBits = caches.AddressMaxNumber - indexBits
)

//DMCacheLine - Direct Mapped Cache Line
type DMCacheLine struct {
	valid bool
	tag   mainMemory.Address
	data  mainMemory.Data
}

type DirectMappedCache struct{
	storage [caches.CacheSize]DMCacheLine
	mainMemory *mainMemory.MainMemory
}

func (dMC *DirectMappedCache)Init(mainMemory *mainMemory.MainMemory){
	dMC.storage = [caches.CacheSize]DMCacheLine{}
	dMC.mainMemory = mainMemory
}

func (dMC *DirectMappedCache) Load(address mainMemory.Address) (mainMemory.Data, bool){
	index, tag:= dMC.extractIndexAndTag(address)
	line :=dMC.storage[index]

	if  line.valid{
		if line.tag == tag {
			return line.data, true
		}
		dMC.mainMemory.Store(line.tag+index,line.data)
	}

	data := dMC.mainMemory.Load(address)

	dMC.storage[index].data = data
	dMC.storage[index].tag = tag
	dMC.storage[index].valid = true

	return data, false
}

func (dMC *DirectMappedCache) Store(address mainMemory.Address, newData mainMemory.Data) bool{
	index, tag := dMC.extractIndexAndTag(address)
	line := dMC.storage[index]

	if line.valid {
		if line.tag == tag {
			dMC.storage[index].data = newData
			return true
		}
		dMC.mainMemory.Store(line.tag+index, line.data)
	}

	dMC.storage[index].data = newData
	dMC.storage[index].tag = tag
	dMC.storage[index].valid = true

	return false
}

func (dMC *DirectMappedCache) extractIndexAndTag(address mainMemory.Address) (mainMemory.Address, mainMemory.Address) {
	index := address & indexBits
	tag := address & tagBits
	return index, tag
}