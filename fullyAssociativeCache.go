package caches

import "sync"

//FACacheLine - Fully Associative Cache Line
type FACacheLine struct {
	valid bool
	address uint32
	data int32
	rWM sync.RWMutex
}

type fullyAssociativeCache struct{
	mainMemory *mainMemory
	storage [cacheSize]FACacheLine
	isStorageFull bool
	lruQueue queue
}

func (fAC *fullyAssociativeCache) Init(mainMemory *mainMemory){
	fAC.mainMemory = mainMemory
	fAC.storage = [cacheSize]FACacheLine{}
	fAC.lruQueue.Init(cacheSize)
}

func (fAC *fullyAssociativeCache) Load(address uint32) (int32, bool){
	line, exist := fAC.getExistingLine(address)
	if exist {
		line.rWM.RUnlock()
		return line.data, exist
	}

	data := fAC.mainMemory.Load(address)

	if !fAC.isStorageFull{
		for index := range fAC.storage {
			line := &fAC.storage[index]

			line.rWM.RLock()
			if !line.valid {
				line.rWM.RUnlock()
				fAC.newAddressInLine(uint32(index), address, data)
				return data, false
			}
			line.rWM.RUnlock()
		}

		fAC.isStorageFull = true
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(indexOfLRU, address, data)

	return data, false
}

func (fAC *fullyAssociativeCache) Store(address uint32, newData int32) bool{
	line, exist := fAC.getExistingLine(address)
	if exist {
		line.rWM.RUnlock()
		line.rWM.Lock()
		line.data = newData
		line.rWM.Unlock()
		return exist
	}

	if !fAC.isStorageFull{
		for index := range fAC.storage {
			line := &fAC.storage[index]

			line.rWM.RLock()
			if !line.valid{
				line.rWM.RUnlock()
				fAC.newAddressInLine(uint32(index), address, newData)
				return false
			}
			line.rWM.RUnlock()
		}
		fAC.isStorageFull = true
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(indexOfLRU, address, newData)

	return false
}

func (fAC *fullyAssociativeCache) getExistingLine(address uint32) (*FACacheLine, bool) {
	for index := range fAC.storage {
		line := &fAC.storage[index]

		line.rWM.RLock()
		if line.address == address && line.valid{
			fAC.lruQueue.Update(uint32(index))
			return line, true
		}
		line.rWM.RUnlock()
	}

	return nil, false
}

func (fAC *fullyAssociativeCache) newAddressInLine(index uint32, address uint32, data int32){
	line := &fAC.storage[index]
	line.rWM.Lock()
	if line.valid{
		oldAddress := line.address
		oldData := line.data
		fAC.mainMemory.Store(oldAddress,oldData)
	}

	fAC.lruQueue.Update(index)
	line.valid = true
	line.address = address
	line.data = data
	line.rWM.Unlock()
}






