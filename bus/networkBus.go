package bus

import (
	"Caches/mainMemory"
	"sync"
)

type LineState string

const(
	NULL = ""
	Modify  = "Modify"
	Exclusive = "Exclusive"
	Shared = "Shared"
	Invalid = "Invalid"
)

type BusMessage struct{
	Address   mainMemory.Address
	LineState LineState
}

type NetworkBus struct {
	numOfListener int
	busListeners [] chan BusMessage
	busWriters [] chan bool
	currentCacheNum int
	busLocker sync.Mutex
}

func (nB *NetworkBus) Init(numOfCaches int) {
	nB.busListeners = make([] chan BusMessage, numOfCaches)
	nB.busWriters = make([] chan bool, numOfCaches)
}

func (nB *NetworkBus) GetBusListenerAndWriter() (chan BusMessage, chan bool) {
	nB.currentCacheNum++
	return nB.busListeners[nB.currentCacheNum], nB.busWriters[nB.currentCacheNum]
}

func (nB *NetworkBus) AskModify(cacheNumber int, address mainMemory.Address) {
	bM := BusMessage{Address: address, LineState: Modify}
	nB.writeAndWaitOnBus(cacheNumber, bM)
}

func (nB *NetworkBus) AskExclusive(cacheNumber int, address mainMemory.Address) bool{
	bM := BusMessage{Address: address, LineState: Exclusive}

	nB.busLocker.Lock()
	for i:=0;i<nB.numOfListener;i++{
		if i!=cacheNumber{
			nB.busListeners[i] <- bM
		}
	}

	var confirmExclusive bool
	for i:=0;i<nB.numOfListener;i++{
		if i!=cacheNumber{
			confirmExclusive = <- nB.busWriters[i]
			if !confirmExclusive{
				i++
				for ;i<nB.numOfListener;i++{
					<- nB.busWriters[i]
				}
				return false
			}
		}
	}
	nB.busLocker.Unlock()

	return true
}

func (nB *NetworkBus) AskShared(cacheNumber int, address mainMemory.Address){
	bM := BusMessage{Address: address, LineState: Shared}
	nB.writeAndWaitOnBus(cacheNumber, bM)
}

func (nB *NetworkBus) writeAndWaitOnBus(cacheNumber int, bM BusMessage) {
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