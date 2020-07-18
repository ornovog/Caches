package N_WayAssociativeCache

import (
	"Caches"
	"Caches/LRU"
	"Caches/mainMemory"
)

const (
	NumOfWays = 2

	//NWayNumBits - the bits of the index in the address
	NWayNumBits = NumOfWays - 1

	//NWayTagBits - the bits of the tag in the address
	NWayTagBits = caches.AddressMaxNumber - NWayNumBits
)

//NWACacheLine - N-Way Associative Cache Line
type NWACacheLine struct {
	valid bool
	tag   mainMemory.Address
	data  mainMemory.Data
}

//NWayAssociativeCache - N-Way Associative Cache
type NWayAssociativeCache struct {
	mainMemory    *mainMemory.MainMemory
	storage       [NumOfWays][caches.CacheSize / NumOfWays]NWACacheLine
	isStorageFull [NumOfWays]bool
	lruQueues     [NumOfWays]LRU.LruQueue
}

//Init - main memory and cache storage Initialization
func (nWAC *NWayAssociativeCache) Init(mainMemory *mainMemory.MainMemory) {
	nWAC.storage = [NumOfWays][caches.CacheSize / NumOfWays]NWACacheLine{}
	nWAC.mainMemory = mainMemory

	for i := 0; i < NumOfWays; i++ {
		nWAC.lruQueues[i].Init(caches.CacheSize / NumOfWays)
	}
}

//Fetch - getting data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Fetch(address mainMemory.Address) (mainMemory.Data, bool) {
	wayNum, tag := extractWayNumAndTag(address)

	line, exist := nWAC.getExistingLine(wayNum, tag)
	if exist {
		return line.data, exist
	}

	data := nWAC.mainMemory.Load(address)

	if !nWAC.isStorageFull[wayNum] {
		for index := range nWAC.storage[wayNum] {
			line = &nWAC.storage[wayNum][index]
			if !line.valid {
				nWAC.newAddressInLine(wayNum, mainMemory.Address(index), tag, data)
				return data, false
			}
		}
		nWAC.isStorageFull[wayNum] = true
	}

	indexOfLRU := nWAC.lruQueues[wayNum].Back()
	nWAC.newAddressInLine(wayNum, mainMemory.Address(indexOfLRU), tag, data)

	return data, false
}

//Store - updating data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Store(address mainMemory.Address, newData mainMemory.Data) bool {
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
				nWAC.newAddressInLine(wayNum, mainMemory.Address(index), tag, newData)
				return false
			}
		}
		nWAC.isStorageFull[wayNum] = true
	}

	indexOfLRU := nWAC.lruQueues[wayNum].Back()
	nWAC.newAddressInLine(wayNum, mainMemory.Address(indexOfLRU), tag, newData)

	return false
}

func extractWayNumAndTag(address mainMemory.Address) (mainMemory.Address, mainMemory.Address) {
	wayNum := address & NWayNumBits
	tag := address & NWayTagBits
	return wayNum, tag
}

func (nWAC *NWayAssociativeCache) getExistingLine(wayNum, tag mainMemory.Address) (*NWACacheLine, bool) {
	for index := range nWAC.storage[wayNum] {
		line := &nWAC.storage[wayNum][index]

		if line.tag == tag && line.valid {
			nWAC.lruQueues[wayNum].Update(mainMemory.Address(index))
			return line, true
		}
	}

	return nil, false
}

func (nWAC *NWayAssociativeCache) newAddressInLine(wayNum, index, tag mainMemory.Address, data mainMemory.Data) {
	line := &nWAC.storage[wayNum][index]

	if line.valid {
		oldAddress := line.tag + wayNum
		oldData := line.data
		nWAC.mainMemory.Store(oldAddress, oldData)
	}

	nWAC.lruQueues[wayNum].Update(index)
	line.valid = true
	line.tag = tag
	line.data = data
}
