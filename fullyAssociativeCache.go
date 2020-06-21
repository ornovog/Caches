package Caches

type FACacheLine struct {
	useNumber uint64
	address uint32
	data float64
}

type fullyAssociativeCache struct{
	useNumber uint64
	storage [cacheSize]FACacheLine
	isStorageFull bool
	mM *mainMemory
}

func (fAC *fullyAssociativeCache) GetData(address uint32) (float64, bool){
	data, exist := fAC.getExistingLine(address)
	if exist {
		return data, exist
	}

	data = fAC.mM.Fetch(address)

	if !fAC.isStorageFull{
		for i, line := range fAC.storage {
			if line.useNumber == 0 {
				fAC.updateInIndex(uint32(i), address, data)
				return data, false
			}
		}
	}else {
		indexOfLRU := fAC.lRU()
		fAC.updateInIndex(indexOfLRU, address, data)
	}
	return data, false
}

func (fAC *fullyAssociativeCache) getExistingLine(address uint32) (float64, bool) {
	for _, line := range fAC.storage {
		if line.address == address {
			line.useNumber = fAC.newUseNumber()
			return line.data, true
		}
	}

	return 0, false
}

func (fAC *fullyAssociativeCache) newUseNumber()uint64{
	fAC.useNumber++
	return fAC.useNumber
}

func (fAC *fullyAssociativeCache) updateInIndex(index uint32, address uint32, data float64){
	fAC.storage[index].useNumber = fAC.newUseNumber()
	fAC.storage[index].address = address
	fAC.storage[index].data = data
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








