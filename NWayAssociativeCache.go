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
	useNumber uint64
	storage [numOfWays][cacheSize/numOfWays]NWACacheLine
	isStorageFull [numOfWays]bool
	mM *mainMemory
}

func (nWAC *NWayAssociativeCache) GetData(address uint32) (byte, bool){
	wayIndex, tag := extractWayIndexAndTag(address)

	line, exist := nWAC.getExistingLine(wayIndex,tag)
	if exist {
		return line.data, exist
	}

	data := nWAC.mM.Fetch(address)

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
		if line.tag == tag {
			line.useNumber = nWAC.newUseNumber()
			return &line, true
		}
	}

	return 0, false
}

func (nWAC *NWayAssociativeCache) newUseNumber()uint64{
	nWAC.useNumber++
	return nWAC.useNumber
}

func (nWAC *NWayAssociativeCache) newAddressInLine(wayIndex, index, tag uint32, data byte){
	line := &nWAC.storage[wayIndex][index]
	oldAddress := line.tag + wayIndex
	oldData := line.data
	nWAC.mM.Store(oldAddress,oldData)

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
