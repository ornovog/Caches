package bus

import (
	"Caches/mainMemory"
	"sync"
)

type LineState string

const (
	NULL      = ""
	Modify    = "Modify"
	Exclusive = "Exclusive"
	Shared    = "Shared"
	Invalid   = "Invalid"
)

type Message struct {
	Address   mainMemory.Address
	LineState LineState
}

type NetworkBus struct {
	numOfListener   int
	busListeners    []chan Message
	busWriters      []chan bool
	currentCacheNum int
	busLocker       sync.Mutex
}

func (nB *NetworkBus) Init(numOfCaches int) {
	nB.busListeners = make([]chan Message, numOfCaches)
	nB.busWriters = make([]chan bool, numOfCaches)
}

func (nB *NetworkBus) GetBusListenerAndWriter() (chan Message, chan bool, int) {
	cacheNum := nB.currentCacheNum
	nB.currentCacheNum++
	return nB.busListeners[cacheNum], nB.busWriters[cacheNum], cacheNum
}

func (nB *NetworkBus) AskModify(cacheNumber int, address mainMemory.Address) {
	bM := Message{Address: address, LineState: Modify}
	nB.writeAndWaitOnBus(cacheNumber, bM)
}

func (nB *NetworkBus) AskExclusive(cacheNumber int, address mainMemory.Address) bool {
	bM := Message{Address: address, LineState: Exclusive}

	nB.busLocker.Lock()
	for i := 0; i < nB.numOfListener; i++ {
		if i != cacheNumber {
			nB.busListeners[i] <- bM
		}
	}

	var confirmExclusive bool
	for i := 0; i < nB.numOfListener; i++ {
		if i != cacheNumber {
			confirmExclusive = <-nB.busWriters[i]
			if !confirmExclusive {
				i++
				for ; i < nB.numOfListener; i++ {
					<-nB.busWriters[i]
				}
				return false
			}
		}
	}
	nB.busLocker.Unlock()

	return true
}

func (nB *NetworkBus) AskShared(cacheNumber int, address mainMemory.Address) {
	bM := Message{Address: address, LineState: Shared}
	nB.writeAndWaitOnBus(cacheNumber, bM)
}

func (nB *NetworkBus) writeAndWaitOnBus(cacheNumber int, bM Message) {
	nB.busLocker.Lock()
	for i := 0; i < nB.numOfListener; i++ {
		if i != cacheNumber {
			nB.busListeners[i] <- bM
		}
	}

	for i := 0; i < nB.numOfListener; i++ {
		if i != cacheNumber {
			_ = <-nB.busWriters[i]
		}
	}
	nB.busLocker.Unlock()
}
