package Caches

const (
	numOfWays = 2
	NWayIndexBits = numOfWays-1
	NWayTagBits = addressMaxNumber-NWayIndexBits
)

type NWACacheLine struct {
	useNumber uint64
	tag uint32
	data float64
}

type NWayAssociativeCache struct{
	useNumber uint64
	storage [numOfWays][cacheSize/numOfWays]NWACacheLine
	isStorageFull [numOfWays]bool
	mM *mainMemory
}

func (nWAC *NWayAssociativeCache) GetData(address uint32) (float64, bool){
	wayIndex := address & NWayIndexBits
	tag := address & NWayTagBits

	data, exist := nWAC.getExistingLine(wayIndex,tag)
	if exist {
		return data, exist
	}

	data = nWAC.mM.Fetch(address)

	if !nWAC.isStorageFull[wayIndex]{
		for index, line := range nWAC.storage[wayIndex] {
			if line.useNumber == 0 {
				nWAC.updateInIndex(wayIndex, uint32(index), tag, data)
				return data, false
			}
		}
	}else {
		indexOfLRU := nWAC.lRU(wayIndex)
		nWAC.updateInIndex(wayIndex, indexOfLRU, tag, data)
	}
	return data, false
}

func (nWAC *NWayAssociativeCache) getExistingLine(wayNum, tag uint32) (float64, bool) {
	for _, line := range nWAC.storage[wayNum] {
		if line.tag == tag {
			line.useNumber = nWAC.newUseNumber()
			return line.data, true
		}
	}

	return 0, false
}

func (nWAC *NWayAssociativeCache) newUseNumber()uint64{
	nWAC.useNumber++
	return nWAC.useNumber
}

func (nWAC *NWayAssociativeCache) updateInIndex(wayNum, index, tag uint32, data float64){
	nWAC.storage[wayNum][index].useNumber = nWAC.newUseNumber()
	nWAC.storage[wayNum][index].tag = tag
	nWAC.storage[wayNum][index].data = data
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