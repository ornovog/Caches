package Caches

const (
	numOfWays = 2
	NWayIndexBits = numOfWays-1
	NWayTagBits = addressMaxNumber-NWayIndexBits
)

type NWACacheLine struct {
	useNumber uint64
	tag uint32
	data byte
}

type NWayAssociativeCache struct{
	useNumberCounter uint64
	storage *[numOfWays][cacheSize/numOfWays]NWACacheLine
	isStorageFull [numOfWays]bool
	mainMemory *mainMemory
}

func (nWAC *NWayAssociativeCache)Init(mainMemory *mainMemory){
	nWAC.storage = &[numOfWays][cacheSize/numOfWays]NWACacheLine{}
	nWAC.mainMemory = mainMemory
}

func (nWAC *NWayAssociativeCache) GetData(address uint32) (byte, bool){
	wayIndex, tag := extractWayIndexAndTag(address)

	line, exist := nWAC.getExistingLine(wayIndex,tag)
	if exist {
		return line.data, exist
	}

	data := nWAC.mainMemory.Fetch(address)

	if !nWAC.isStorageFull[wayIndex]{
		for index, line := range nWAC.storage[wayIndex] {
			if line.useNumber == 0 {
				nWAC.newAddressInLine(wayIndex, uint32(index), tag, data)
				return data, false
			}
		}
	}else {
		indexOfLRU := nWAC.lRU(wayIndex)
		nWAC.newAddressInLine(wayIndex, indexOfLRU, tag, data)
	}

	return data, false
}

func extractWayIndexAndTag(address uint32) (uint32, uint32) {
	wayIndex := address & NWayIndexBits
	tag := address & NWayTagBits
	return wayIndex, tag
}

func (nWAC *NWayAssociativeCache) getExistingLine(wayNum, tag uint32) (*NWACacheLine, bool) {
	for _, line := range nWAC.storage[wayNum] {
		if line.tag == tag && line.useNumber!=0{
			line.useNumber = nWAC.newUseNumber()
			return &line, true
		}
	}

	return nil, false
}

func (nWAC *NWayAssociativeCache) newUseNumber()uint64{
	nWAC.useNumberCounter++
	return nWAC.useNumberCounter
}

func (nWAC *NWayAssociativeCache) newAddressInLine(wayIndex, index, tag uint32, data byte){
	line := &nWAC.storage[wayIndex][index]
	if line.useNumber != 0{
		oldAddress := line.tag + wayIndex
		oldData := line.data
		nWAC.mainMemory.Store(oldAddress,oldData)
	}

	line.useNumber = nWAC.newUseNumber()
	line.tag = tag
	line.data = data
}

func (nWAC *NWayAssociativeCache) lRU(wayNum uint32) uint32 {
	indexOfLRU := 0
	minUseNumber := nWAC.storage[wayNum][0].useNumber

	for i, line := range nWAC.storage[wayNum] {
		if line.useNumber < minUseNumber {
			indexOfLRU = i
			minUseNumber = line.useNumber
		}
	}
	return uint32(indexOfLRU)
}

func (nWAC *NWayAssociativeCache) Update(address uint32, newData byte) bool{
	wayIndex, tag := extractWayIndexAndTag(address)

	line, exist := nWAC.getExistingLine(wayIndex, tag)
	if exist {
		line.data = newData
		return exist
	}

	if !nWAC.isStorageFull[wayIndex] {
		for index, line := range nWAC.storage[wayIndex] {
			if line.useNumber == 0 {
				nWAC.newAddressInLine(wayIndex, uint32(index), tag, newData)
				return false
			}
		}
	} else {
		indexOfLRU := nWAC.lRU(wayIndex)
		nWAC.newAddressInLine(wayIndex, indexOfLRU, tag, newData)
	}

	return false
}
