package caches

const (
	numOfWays = 2

	//NWayIndexBits - the bits of the index in the address
	NWayIndexBits = numOfWays-1

	//NWayTagBits - the bits of the tag in the address
	NWayTagBits = addressMaxNumber-NWayIndexBits
)

//NWACacheLine - N-Way Associative Cache Line
type NWACacheLine struct {
	useNumber uint64
	tag uint32
	data byte
}

//NWayAssociativeCache - N-Way Associative Cache
type NWayAssociativeCache struct{
	useNumberCounter uint64
	storage *[numOfWays][cacheSize/numOfWays]NWACacheLine
	isStorageFull [numOfWays]bool
	mainMemory *mainMemory
}

//Init - main memory and cache storage Initialization
func (nWAC *NWayAssociativeCache)Init(mainMemory *mainMemory){
	nWAC.storage = &[numOfWays][cacheSize/numOfWays]NWACacheLine{}
	nWAC.mainMemory = mainMemory
}

//Fetch - getting data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Fetch(address uint32) (byte, bool){
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
		nWAC.isStorageFull[wayIndex] = true
	}


	indexOfLRU := nWAC.lRU(wayIndex)
	nWAC.newAddressInLine(wayIndex, indexOfLRU, tag, data)

	return data, false
}

//Store - updating data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Store(address uint32, newData byte) bool{
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
		nWAC.isStorageFull[wayIndex] = true
	}

	indexOfLRU := nWAC.lRU(wayIndex)
	nWAC.newAddressInLine(wayIndex, indexOfLRU, tag, newData)

	return false
}

func extractWayIndexAndTag(address uint32) (uint32, uint32) {
	wayIndex := address & NWayIndexBits
	tag := address & NWayTagBits
	return wayIndex, tag
}

func (nWAC *NWayAssociativeCache) getExistingLine(wayNum, tag uint32) (*NWACacheLine, bool) {
	for i, line := range nWAC.storage[wayNum] {
		if line.tag == tag && line.useNumber!=0{
			nWAC.storage[wayNum][i].useNumber = nWAC.newUseNumber()
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
