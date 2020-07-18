package fullyAssociativeCache

import (
	"Caches"
	"Caches/LRU"
	"Caches/bus"
	"Caches/mainMemory"
)

//FACacheLine - Fully Associative Cache Line
type FACacheLine struct {
	state   bus.LineState
	address mainMemory.Address
	data    mainMemory.Data
}

type fullyAssociativeCache struct {
	mainMemory     *mainMemory.MainMemory
	networkBus     *bus.NetworkBus
	cacheNumber    int
	busListener    chan bus.Message
	busWriter      chan bool
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
	fAC.networkBus = networkBus
	fAC.busListener, fAC.busWriter, fAC.cacheNumber = networkBus.GetBusListenerAndWriter()
	go fAC.listenOnBus()
}

func (fAC *fullyAssociativeCache) listenOnBus() {
	for busMessage := range fAC.busListener {
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
			fAC.busWriter <- true
		}
	}
}

//MESI functions
func (fAC *fullyAssociativeCache) sharedCase(index int, busMessage bus.Message) {
	if fAC.storage[index].state == bus.Modify || fAC.storage[index].state == bus.Exclusive {
		fAC.mainMemory.Store(busMessage.Address, fAC.storage[index].data)
		fAC.storage[index].state = bus.Shared
	}
	fAC.busWriter <- true
}

func (fAC *fullyAssociativeCache) exclusiveCase(index int) {
	if fAC.storage[index].state == bus.NULL {
		fAC.busWriter <- true
	} else {
		fAC.busWriter <- false
	}
}

func (fAC *fullyAssociativeCache) modifyCase(index int, busMessage bus.Message) {
	if fAC.storage[index].state == bus.Modify {
		fAC.mainMemory.Store(busMessage.Address, fAC.storage[index].data)
	}
	fAC.storage[index].state = bus.Invalid
	fAC.busWriter <- true
}

//Load functions
func (fAC *fullyAssociativeCache) Load(address mainMemory.Address) (mainMemory.Data, bool) {
	line, exist := fAC.getExistingLine(address)
	if exist && line.state != bus.Invalid {
		return line.data, exist
	}

	fAC.networkBus.AskShared(fAC.cacheNumber, address)
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
	if exist {
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
	for index := range fAC.storage {
		line := &fAC.storage[index]

		if line.address == address && line.state != bus.NULL {
			fAC.lruQueue.Update(mainMemory.Address(index))
			return line, true
		}
	}

	return nil, false
}

func (fAC *fullyAssociativeCache) newAddressInLine(index mainMemory.Address, address mainMemory.Address, data mainMemory.Data) {
	line := &fAC.storage[index]

	if line.state == bus.Exclusive || line.state == bus.Modify {
		oldAddress := line.address
		oldData := line.data
		fAC.mainMemory.Store(oldAddress, oldData)
	}

	fAC.lruQueue.Update(index)
	line.state = bus.Shared
	line.address = address
	line.data = data
}
