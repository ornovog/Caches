package Caches

const (
	numOfWays = 2
	wayNumberBits = 1
	twoWayTagBits = 4294967294
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
	wayNum := address & wayNumberBits
	tag := address & twoWayTagBits

	data, exist := nWAC.getExistingLine(wayNum,tag)
	if exist {
		return data, exist
	}

	data = nWAC.mM.Fetch(address)

	if !nWAC.isStorageFull[wayNum]{
		for index, line := range nWAC.storage[wayNum] {
			if line.useNumber == 0 {
				nWAC.updateInIndex(wayNum, uint32(index), tag, data)
				return data, false
			}
		}
	}else {
		indexOfLRU := nWAC.lRU(wayNum)
		nWAC.updateInIndex(wayNum, indexOfLRU, tag, data)
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