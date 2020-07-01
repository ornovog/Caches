package caches

import (
	"sync"
)

const (
	indexBits = cacheSize - 1
	tagBits = addressMaxNumber - indexBits
)

//DMCacheLine - Direct Mapped Cache Line
type DMCacheLine struct {
	valid bool
	tag uint32
	data int32
	rWM sync.RWMutex
}

type directMappedCache struct{
	storage [cacheSize]DMCacheLine
	mainMemory *mainMemory
}

func (dMC *directMappedCache)Init(mainMemory *mainMemory){
	dMC.storage = [cacheSize]DMCacheLine{}
	dMC.mainMemory = mainMemory
}

func (dMC *directMappedCache) Load(address uint32) (int32, bool){
	index, tag:= dMC.extractIndexAndTag(address)
	line :=dMC.storage[index]

	line.rWM.RLock()
	if  line.valid{
		if line.tag == tag {
			line.rWM.RUnlock()
			return line.data, true
		}
		dMC.mainMemory.Store(line.tag+index,line.data)
	}
	line.rWM.RUnlock()

	data := dMC.mainMemory.Load(address)

	line.rWM.Lock()
	dMC.storage[index].data = data
	dMC.storage[index].tag = tag
	dMC.storage[index].valid = true
	line.rWM.Unlock()

	return data, false
}

func (dMC *directMappedCache) Store(address uint32, newData int32) bool{
	index, tag := dMC.extractIndexAndTag(address)
	line := dMC.storage[index]

	line.rWM.RLock()
	if line.valid {
		if line.tag == tag {
			dMC.storage[index].data = newData
			line.rWM.RUnlock()
			return true
		}
		dMC.mainMemory.Store(line.tag+index, line.data)
	}
	line.rWM.RUnlock()

	line.rWM.Lock()
	dMC.storage[index].data = newData
	dMC.storage[index].tag = tag
	dMC.storage[index].valid = true
	line.rWM.Unlock()

	return false
}

func (dMC *directMappedCache) extractIndexAndTag(address uint32) (uint32, uint32) {
	index := address & indexBits
	tag := address & tagBits
	return index, tag
}