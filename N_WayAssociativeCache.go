package caches

const (
	numOfWays = 2

	//NWayNumBits - the bits of the index in the address
	NWayNumBits = numOfWays-1

	//NWayTagBits - the bits of the tag in the address
	NWayTagBits = addressMaxNumber-NWayNumBits
)

//NWACacheLine - N-Way Associative Cache Line
type NWACacheLine struct {
	valid bool
	tag uint32
	data byte
}

//NWayAssociativeCache - N-Way Associative Cache
type NWayAssociativeCache struct{
	mainMemory *mainMemory
	storage [numOfWays][cacheSize/numOfWays]NWACacheLine
	isStorageFull [numOfWays]bool
	lruQueues [numOfWays]queue
}

//Init - main memory and cache storage Initialization
func (nWAC *NWayAssociativeCache)Init(mainMemory *mainMemory){
	nWAC.storage = [numOfWays][cacheSize/numOfWays]NWACacheLine{}
	nWAC.mainMemory = mainMemory

	for i :=0; i< numOfWays; i++ {
		nWAC.lruQueues[i].Init(cacheSize/numOfWays)
	}
}

//Fetch - getting data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Fetch(address uint32) (byte, bool){
	wayNum, tag := extractwayNumAndTag(address)


	line, exist := nWAC.getExistingLine(wayNum,tag)
	if exist {
		return line.data, exist
	}

	data := nWAC.mainMemory.Fetch(address)

	if !nWAC.isStorageFull[wayNum]{
		for index, line := range nWAC.storage[wayNum] {
			if !line.valid {
				nWAC.newAddressInLine(wayNum, uint32(index), tag, data)
				return data, false
			}
		}
		nWAC.isStorageFull[wayNum] = true
	}


	indexOfLRU := nWAC.lruQueues[wayNum].Back()
	nWAC.newAddressInLine(wayNum, indexOfLRU, tag, data)

	return data, false
}

//Store - updating data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Store(address uint32, newData byte) bool{
	wayNum, tag := extractwayNumAndTag(address)

	line, exist := nWAC.getExistingLine(wayNum, tag)
	if exist {
		line.data = newData
		return exist
	}

	if !nWAC.isStorageFull[wayNum] {
		for index, line := range nWAC.storage[wayNum] {
			if !line.valid {
				nWAC.newAddressInLine(wayNum, uint32(index), tag, newData)
				return false
			}
		}
		nWAC.isStorageFull[wayNum] = true
	}

	indexOfLRU := nWAC.lruQueues[wayNum].Back()
	nWAC.newAddressInLine(wayNum, indexOfLRU, tag, newData)

	return false
}

func extractwayNumAndTag(address uint32) (uint32, uint32) {
	wayNum := address & NWayNumBits
	tag := address & NWayTagBits
	return wayNum, tag
}

func (nWAC *NWayAssociativeCache) getExistingLine(wayNum, tag uint32) (*NWACacheLine, bool) {
	for index, line := range nWAC.storage[wayNum] {
		if line.tag == tag && line.valid{
			nWAC.lruQueues[wayNum].Update(uint32(index))
			return &line, true
		}
	}

	return nil, false
}

func (nWAC *NWayAssociativeCache) newAddressInLine(wayNum, index, tag uint32, data byte){
	line := &nWAC.storage[wayNum][index]
	if line.valid {
		oldAddress := line.tag + wayNum
		oldData := line.data
		nWAC.mainMemory.Store(oldAddress,oldData)
	}

	nWAC.lruQueues[wayNum].Update(index)
	line.valid = true
	line.tag = tag
	line.data = data
}
