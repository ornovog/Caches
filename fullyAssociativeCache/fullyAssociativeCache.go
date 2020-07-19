package fullyAssociativeCache

import (
	"Caches"
	"Caches/LRU"
	"Caches/bus"
	"Caches/mainMemory"
	"sync"
)

//FACacheLine - Fully Associative Cache Line
type FACacheLine struct {
	state   	bus.LineState
	address 	mainMemory.Address
	data    	mainMemory.Data
	stateLocker sync.Mutex
}

type fullyAssociativeCache struct {
	mainMemory     *mainMemory.MainMemory
	NetworkBus     *bus.NetworkBus
	busStuff       bus.BusStuff
	storage        [caches.CacheSize]FACacheLine
	nextEmpty      int
	lruQueue       LRU.LruQueue
	addressToIndex map[mainMemory.Address]int
}

//Init functions
func (fAC *fullyAssociativeCache) Init(mainMemory *mainMemory.MainMemory, networkBus *bus.NetworkBus) {
	fAC.mainMemory = mainMemory
	fAC.storage = [caches.CacheSize]FACacheLine{}
	fAC.lruQueue.Init(caches.CacheSize)
	fAC.busStuff.NetworkBus = networkBus
	fAC.busStuff = networkBus.GetBusListenerAndWriter()
	go fAC.listenOnBus()
}

func (fAC *fullyAssociativeCache) listenOnBus() {
	for busMessage := range fAC.busStuff.BusListener {
		index, ok := fAC.addressToIndex[busMessage.Address]
		if ok {
			switch busMessage.LineState {
			case bus.Modify:
				fAC.modifyCase(index, busMessage)
			case bus.Exclusive:
				fAC.exclusiveCase(index)
			case bus.Shared:
				fAC.sharedCase(index, busMessage)
			}
		} else {
			fAC.busStuff.BusWriter <- true
		}
	}
}

//MESI functions
func (fAC *fullyAssociativeCache) sharedCase(index int, busMessage bus.Message) {
	fAC.storage[index].stateLocker.Lock()
	defer fAC.storage[index].stateLocker.Unlock()

	if fAC.storage[index].state == bus.Modify || fAC.storage[index].state == bus.Exclusive {
		fAC.mainMemory.Store(busMessage.Address, fAC.storage[index].data)
		fAC.storage[index].state = bus.Shared
	}
	fAC.busStuff.BusWriter <- true
}

func (fAC *fullyAssociativeCache) exclusiveCase(index int) {
	fAC.storage[index].stateLocker.Lock()
	defer fAC.storage[index].stateLocker.Unlock()

	if fAC.storage[index].state == bus.NULL {
		fAC.busStuff.BusWriter <- true
	} else {
		fAC.busStuff.BusWriter <- false
	}
}

func (fAC *fullyAssociativeCache) modifyCase(index int, busMessage bus.Message) {
	fAC.storage[index].stateLocker.Lock()
	defer fAC.storage[index].stateLocker.Unlock()

	if fAC.storage[index].state == bus.Modify {
		fAC.mainMemory.Store(busMessage.Address, fAC.storage[index].data)
	}
	fAC.storage[index].state = bus.Invalid
	fAC.busStuff.BusWriter <- true
}

//Load functions
func (fAC *fullyAssociativeCache) Load(address mainMemory.Address) (mainMemory.Data, bool) {

	line, exist := fAC.getExistingLine(address)
	if exist {
		line.stateLocker.Lock()
		defer line.stateLocker.Unlock()
		if line.state == bus.Invalid {
			fAC.NetworkBus.AskShared(fAC.busStuff.CacheNumber, address)
			line.state = bus.Shared
			data := fAC.mainMemory.Load(address)
			line.data = data
			return line.data, false
		}
		return line.data, true
	}

	line.stateLocker.Lock()
	defer line.stateLocker.Unlock()

	fAC.networkBus.AskShared(fAC.cacheNumber, address)
	line.state = bus.Shared
	data := fAC.mainMemory.Load(address)
	if fAC.nextEmpty < len(fAC.storage) {
		fAC.newAddressInLine(mainMemory.Address(fAC.nextEmpty), address, data)
		fAC.nextEmpty++
		return data, false
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(indexOfLRU, address, data)
	return data, false
}

//Store functions
func (fAC *fullyAssociativeCache) Store(address mainMemory.Address, newData mainMemory.Data) bool {
	line, exist := fAC.getExistingLine(address)

	if fAC.NetworkBus.AskExclusive(fAC.busStuff.CacheNumber, address){
		line.stateLocker.Lock()
		line.state = bus.Exclusive
		line.stateLocker.Unlock()
	}
	if exist && line.state == bus.Exclusive {
		line.data = newData
		return exist
	}

	if fAC.nextEmpty < len(fAC.storage) {
		fAC.newAddressInLine(mainMemory.Address(fAC.nextEmpty), address, newData)
		fAC.nextEmpty++
		return false
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(indexOfLRU, address, newData)

	return false
}

func (fAC *fullyAssociativeCache) getExistingLine(address mainMemory.Address) (*FACacheLine, bool) {
	index, ok := fAC.addressToIndex[address]
	if ok {
		fAC.lruQueue.Update(mainMemory.Address(index))
		return line, true
	}

	return nil, false
}

func (fAC *fullyAssociativeCache) newAddressInLine(index mainMemory.Address, address mainMemory.Address, data mainMemory.Data) {
	line := &fAC.storage[index]
	if line.state == bus.Exclusive || line.state == bus.Modify {
		oldAddress := line.address
		oldData := line.data
		fAC.mainMemory.Store(oldAddress, oldData)
		delete(fAC.addressToIndex, oldAddress)

	}

	fAC.lruQueue.Update(index)
	line.address = address
	line.data = data
}
