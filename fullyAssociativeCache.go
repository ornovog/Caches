package Caches

type FACacheLine struct {
	useNumber uint64
	address uint32
	data byte
}

type fullyAssociativeCache struct{
	useNumber uint64
	storage [cacheSize]FACacheLine
	isStorageFull bool
	mM *mainMemory
}

func (fAC *fullyAssociativeCache) GetData(address uint32) (byte, bool){
	line, exist := fAC.getExistingLine(address)
	if exist {
		return line.data, exist
	}

	data := fAC.mM.Fetch(address)

	if !fAC.isStorageFull{
		for i, line := range fAC.storage {
			if line.useNumber == 0 {
				fAC.newAddressInLine(uint32(i), address, data)
				return data, false
			}
		}
	}else {
		indexOfLRU := fAC.lRU()
		fAC.newAddressInLine(indexOfLRU, address, data)
	}
	return data, false
}

func (fAC *fullyAssociativeCache) getExistingLine(address uint32) (*FACacheLine, bool) {
	for _, line := range fAC.storage {
		if line.address == address {
			line.useNumber = fAC.newUseNumber()
			return &line, true
		}
	}

	return nil, false
}

func (fAC *fullyAssociativeCache) newUseNumber()uint64{
	fAC.useNumber++
	return fAC.useNumber
}

func (fAC *fullyAssociativeCache) newAddressInLine(index uint32, address uint32, data byte){
	line := &fAC.storage[index]
	oldAddress := line.address
	oldData := line.data
	fAC.mM.Store(oldAddress,oldData)

	line.useNumber = fAC.newUseNumber()
	line.address = address
	line.data = data
}

func (fAC *fullyAssociativeCache) lRU() uint32 {
	indexOfLRU := 0
	minUseNumber := fAC.storage[0].useNumber

	for i, line := range fAC.storage {
		if line.useNumber < minUseNumber {
			indexOfLRU = i
			minUseNumber = line.useNumber
		}
	}
	return uint32(indexOfLRU)
}

func (fAC *fullyAssociativeCache) Update(address uint32, newData byte) bool{
	line, exist := fAC.getExistingLine(address)
	if exist {
		line.data = newData
		return exist
	}

	if !fAC.isStorageFull{
		for i, line := range fAC.storage {
			if line.useNumber == 0 {
				fAC.newAddressInLine(uint32(i), address, newData)
				return false
			}
		}
	}else {
		indexOfLRU := fAC.lRU()
		fAC.newAddressInLine(indexOfLRU, address, newData)
	}
	return false
}






