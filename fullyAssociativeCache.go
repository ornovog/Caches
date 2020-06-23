package caches

//FACacheLine - Fully Associative Cache Line
type FACacheLine struct {
	useNumber uint64
	address uint32
	data byte
}

type fullyAssociativeCache struct{
	useNumberCounter uint64
	storage *[cacheSize]FACacheLine
	isStorageFull bool
	mainMemory *mainMemory
}

func (fAC *fullyAssociativeCache) Init(mainMemory *mainMemory){
	fAC.storage = &[cacheSize]FACacheLine{}
	fAC.mainMemory = mainMemory
}

func (fAC *fullyAssociativeCache) Fetch(address uint32) (byte, bool){
	line, exist := fAC.getExistingLine(address)
	if exist {
		return line.data, exist
	}

	data := fAC.mainMemory.Fetch(address)

	if !fAC.isStorageFull{
		for i, line := range fAC.storage {
			if line.useNumber == 0 {
				fAC.newAddressInLine(uint32(i), address, data)
				return data, false
			}
		}

		fAC.isStorageFull = true
	}

	indexOfLRU := fAC.lRU()
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
			if line.useNumber == 0 {
				fAC.newAddressInLine(uint32(i), address, newData)
				return false
			}
		}
		fAC.isStorageFull = true
	}

	indexOfLRU := fAC.lRU()
	fAC.newAddressInLine(indexOfLRU, address, newData)

	return false
}

func (fAC *fullyAssociativeCache) getExistingLine(address uint32) (*FACacheLine, bool) {
	for i, line := range fAC.storage {
		if line.address == address && line.useNumber!=0{
			fAC.storage[i].useNumber = fAC.newUseNumber()
			return &fAC.storage[i], true
		}
	}

	return nil, false
}

func (fAC *fullyAssociativeCache) newUseNumber()uint64{
	fAC.useNumberCounter++
	return fAC.useNumberCounter
}

func (fAC *fullyAssociativeCache) newAddressInLine(index uint32, address uint32, data byte){
	line := &fAC.storage[index]
	if line.useNumber !=0{
		oldAddress := line.address
		oldData := line.data
		fAC.mainMemory.Store(oldAddress,oldData)
	}

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






