package caches

import "sync"

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
	data int32
	rWM sync.RWMutex
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
func (nWAC *NWayAssociativeCache) Fetch(address uint32) (int32, bool){
	wayNum, tag := extractwayNumAndTag(address)

	line, exist := nWAC.getExistingLine(wayNum,tag)
	if exist {
		return line.data, exist
	}

	data := nWAC.mainMemory.Load(address)

	if !nWAC.isStorageFull[wayNum]{
		for index := range nWAC.storage[wayNum] {
			line = &nWAC.storage[wayNum][index]
			line.rWM.RLock()
			if !line.valid {
				line.rWM.RUnlock()
				nWAC.newAddressInLine(wayNum, uint32(index), tag, data)
				return data, false
			}
			line.rWM.RUnlock()
		}
		nWAC.isStorageFull[wayNum] = true
	}


	indexOfLRU := nWAC.lruQueues[wayNum].Back()
	nWAC.newAddressInLine(wayNum, indexOfLRU, tag, data)

	return data, false
}

//Store - updating data and if was a hit, by passing address
func (nWAC *NWayAssociativeCache) Store(address uint32, newData int32) bool{
	wayNum, tag := extractwayNumAndTag(address)

	line, exist := nWAC.getExistingLine(wayNum, tag)
	if exist {
		line.rWM.RUnlock()

		line.rWM.RLock()
		line.data = newData
		line.rWM.RUnlock()
		return exist
	}

	if !nWAC.isStorageFull[wayNum] {
		for index := range nWAC.storage[wayNum] {
			line := &nWAC.storage[wayNum][index]

			line.rWM.RLock()
			if !line.valid {
				line.rWM.RUnlock()
				nWAC.newAddressInLine(wayNum, uint32(index), tag, newData)
				return false
			}
			line.rWM.RUnlock()
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
	for index := range nWAC.storage[wayNum] {
		line := &nWAC.storage[wayNum][index]

		line.rWM.RLock()
		if line.tag == tag && line.valid{
			nWAC.lruQueues[wayNum].Update(uint32(index))
			return line, true
		}
		line.rWM.RUnlock()
	}

	return nil, false
}

func (nWAC *NWayAssociativeCache) newAddressInLine(wayNum, index, tag uint32, data int32){
	line := &nWAC.storage[wayNum][index]

	line.rWM.Lock()
	if line.valid {
		oldAddress := line.tag + wayNum
		oldData := line.data
		nWAC.mainMemory.Store(oldAddress,oldData)
	}

	nWAC.lruQueues[wayNum].Update(index)
	line.valid = true
	line.tag = tag
	line.data = data
	line.rWM.Unlock()
}
