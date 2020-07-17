package caches

const (
	NumOfWays = 2

	//NWayNumBits - the bits of the index in the address
	NWayNumBits = NumOfWays -1

	//NWayTagBits - the bits of the tag in the address
	NWayTagBits = AddressMaxNumber -NWayNumBits
)

//NWACacheLine - N-Way Associative Cache Line
type NWACacheLine struct {
	valid bool
	tag Address
	data Data

}

//NWayAssociativeCache - N-Way Associative Cache
type NWayAssociativeCache struct{
	mainMemory *mainMemory
	storage [NumOfWays][CacheSize/NumOfWays]NWACacheLine
	isStorageFull [NumOfWays]bool
	lruQueues [NumOfWays]lruQueue
}

//Init - main memory and cache storage Initialization
func (nWAC *NWayAssociativeCache)Init(mainMemory *mainMemory){
	nWAC.storage = [NumOfWays][CacheSize / NumOfWays]NWACacheLine{}
	nWAC.mainMemory = mainMemory

	for i :=0; i< NumOfWays; i++ {
		nWAC.lruQueues[i].Init(CacheSize / NumOfWays)
	}
}

//Fetch - getting data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Fetch(address Address) (Data, bool){
	wayNum, tag := extractWayNumAndTag(address)

	line, exist := nWAC.getExistingLine(wayNum,tag)
	if exist {
		return line.data, exist
	}

	data := nWAC.mainMemory.Load(address)

	if !nWAC.isStorageFull[wayNum]{
		for index := range nWAC.storage[wayNum] {
			line = &nWAC.storage[wayNum][index]
			if !line.valid {
				nWAC.newAddressInLine(wayNum, Address(index), tag, data)
				return data, false
			}
		}
		nWAC.isStorageFull[wayNum] = true
	}


	indexOfLRU := nWAC.lruQueues[wayNum].Back()
	nWAC.newAddressInLine(wayNum, Address(indexOfLRU), tag, data)

	return data, false
}

//Store - updating data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Store(address Address, newData Data) bool{
	wayNum, tag := extractWayNumAndTag(address)

	line, exist := nWAC.getExistingLine(wayNum, tag)
	if exist {
		line.data = newData
		return exist
	}

	if !nWAC.isStorageFull[wayNum] {
		for index := range nWAC.storage[wayNum] {
			line := &nWAC.storage[wayNum][index]

			if !line.valid {
				nWAC.newAddressInLine(wayNum, Address(index), tag, newData)
				return false
			}
		}
		nWAC.isStorageFull[wayNum] = true
	}

	indexOfLRU := nWAC.lruQueues[wayNum].Back()
	nWAC.newAddressInLine(wayNum, Address(indexOfLRU), tag, newData)

	return false
}

func extractWayNumAndTag(address Address) (Address, Address) {
	wayNum := address & NWayNumBits
	tag := address & NWayTagBits
	return wayNum, tag
}

func (nWAC *NWayAssociativeCache) getExistingLine(wayNum, tag Address) (*NWACacheLine, bool) {
	for index := range nWAC.storage[wayNum] {
		line := &nWAC.storage[wayNum][index]

		if line.tag == tag && line.valid{
			nWAC.lruQueues[wayNum].Update(Address(index))
			return line, true
		}
	}

	return nil, false
}

func (nWAC *NWayAssociativeCache) newAddressInLine(wayNum, index, tag Address, data Data){
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
