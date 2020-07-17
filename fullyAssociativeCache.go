package caches

//FACacheLine - Fully Associative Cache Line
type FACacheLine struct {
	state LineState
	address Address
	data Data
}

type fullyAssociativeCache struct{
	mainMemory *mainMemory
	networkBus *networkBus
	cacheNumber int
	busListener chan busMessage
	busWriter chan bool
	storage [CacheSize]FACacheLine
	isStorageFull bool
	lruQueue lruQueue
	addressToIndex map[Address]int
}

//Init functions
func (fAC *fullyAssociativeCache) Init(mainMemory *mainMemory, networkBus *networkBus, cacheNumber int){
	fAC.mainMemory = mainMemory
	fAC.storage = [CacheSize]FACacheLine{}
	fAC.lruQueue.Init(CacheSize)
	fAC.networkBus = networkBus
	fAC.busListener, fAC.busWriter = networkBus.GetBusListenerAndWriter()
	fAC.cacheNumber = cacheNumber
	fAC.listenOnBus()
}

func (fAC *fullyAssociativeCache) listenOnBus(){
	for busMessage := range fAC.busListener {
		index, ok := fAC.addressToIndex[busMessage.Address]
		if ok{
			switch busMessage.LineState {
			case Modify:
				fAC.modifyCase(index, busMessage)
			case Exclusive:
				fAC.exclusiveCase(index)
			case Shared:
				fAC.sharedCase(index, busMessage)
			}
		}else {
			fAC.busWriter <- true
		}
	}
}

//MESI functions
func (fAC *fullyAssociativeCache) sharedCase(index int, busMessage busMessage) {
	if fAC.storage[index].state == Modify || fAC.storage[index].state == Exclusive {
		fAC.mainMemory.Store(busMessage.Address, fAC.storage[index].data)
		fAC.storage[index].state = Shared
	}
	fAC.busWriter <- true
}

func (fAC *fullyAssociativeCache) exclusiveCase(index int) {
	if fAC.storage[index].state == NULL {
		fAC.busWriter <- true
	} else {
		fAC.busWriter <- false
	}
}

func (fAC *fullyAssociativeCache) modifyCase(index int, busMessage busMessage) {
	if fAC.storage[index].state == Modify {
		fAC.mainMemory.Store(busMessage.Address, fAC.storage[index].data)
	}
	fAC.storage[index].state = Invalid
	fAC.busWriter <- true
}

//Load functions
func (fAC *fullyAssociativeCache) Load(address Address) (Data, bool){
	line, exist := fAC.getExistingLine(address)
	if exist && line.state != Invalid {
		return line.data, exist
	}

	fAC.networkBus.AskShared(fAC.cacheNumber, address)
	data := fAC.mainMemory.Load(address)

	if !fAC.isStorageFull{
		found := fAC.findFirstEmpty(address, data)
		if found {
			return data, found
		}
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(Address(indexOfLRU), address, data)

	return data, false
}

func (fAC *fullyAssociativeCache) findFirstEmpty(address Address, data Data) bool {
	for index := range fAC.storage {
		line := &fAC.storage[index]

		if line.state != NULL {
			fAC.newAddressInLine(Address(index), address, data)
			return true
		}
	}

	fAC.isStorageFull = true
	return  false
}

//Store functions
func (fAC *fullyAssociativeCache) Store(address Address, newData Data) bool{
	line, exist := fAC.getExistingLine(address)
	if exist {
		line.data = newData
		return exist
	}

	if !fAC.isStorageFull{
		for index := range fAC.storage {
			line := &fAC.storage[index]

			if line.state == NULL{
				fAC.newAddressInLine(Address(index), address, newData)
				return false
			}
		}
		fAC.isStorageFull = true
	}

	indexOfLRU := fAC.lruQueue.Back()
	fAC.newAddressInLine(Address(indexOfLRU), address, newData)

	return false
}

func (fAC *fullyAssociativeCache) getExistingLine(address Address) (*FACacheLine, bool) {
	for index := range fAC.storage {
		line := &fAC.storage[index]

		if line.address == address && line.state != NULL{
			fAC.lruQueue.Update(Address(index))
			return line, true
		}
	}

	return nil, false
}

func (fAC *fullyAssociativeCache) newAddressInLine(index Address, address Address, data Data){
	line := &fAC.storage[index]

	if line.state == Exclusive || line.state == Modify{
			oldAddress := line.address
			oldData := line.data
			fAC.mainMemory.Store(oldAddress,oldData)
	}

	fAC.lruQueue.Update(index)
	line.state = Shared
	line.address = address
	line.data = data
}






