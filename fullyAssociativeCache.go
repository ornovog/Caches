package caches

//FACacheLine - Fully Associative Cache Line
type FACacheLine struct {
	valid bool
	address uint32
	data byte
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

func (fAC *fullyAssociativeCache) Fetch(address uint32) (byte, bool){
	line, exist := fAC.getExistingLine(address)
	if exist {
		return line.data, exist
	}

	data := fAC.mainMemory.Fetch(address)

	if !fAC.isStorageFull{
		for i, line := range fAC.storage {
			if !line.valid {
				fAC.newAddressInLine(uint32(i), address, data)
				return data, false
			}
		}

		fAC.isStorageFull = true
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(indexOfLRU, address, data)

	return data, false
}

func (fAC *fullyAssociativeCache) Store(address uint32, newData byte) bool{
	line, exist := fAC.getExistingLine(address)
	if exist {
		line.data = newData
		return exist
	}

	if !fAC.isStorageFull{
		for i, line := range fAC.storage {
			if !line.valid{
				fAC.newAddressInLine(uint32(i), address, newData)
				return false
			}
		}
		fAC.isStorageFull = true
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(indexOfLRU, address, newData)

	return false
}

func (fAC *fullyAssociativeCache) getExistingLine(address uint32) (*FACacheLine, bool) {
	for i, line := range fAC.storage {
		if line.address == address && line.valid{
			fAC.lruQueue.Update(uint32(i))
			return &fAC.storage[i], true
		}
	}

	return nil, false
}

func (fAC *fullyAssociativeCache) newAddressInLine(index uint32, address uint32, data byte){
	line := &fAC.storage[index]
	if line.valid{
		oldAddress := line.address
		oldData := line.data
		fAC.mainMemory.Store(oldAddress,oldData)
	}

	fAC.lruQueue.Update(index)
	line.valid = true
	line.address = address
	line.data = data
}






