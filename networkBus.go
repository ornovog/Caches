package caches

import "sync"

type LineState string

const(
	NULL = ""
	Modify  = "Modify"
	Exclusive = "Exclusive"
	Shared = "Shared"
	Invalid = "Invalid"
)

type busMessage struct{
	Address   Address
	LineState LineState
}

type networkBus struct {
	numOfListener int
	busListeners [] chan busMessage
	busWriters [] chan bool
	currentCacheNum int
	busLocker sync.Mutex
}

func (nB *networkBus) Init(numOfCaches int) {
	nB.busListeners = make([] chan busMessage, numOfCaches)
	nB.busWriters = make([] chan bool, numOfCaches)
}

func (nB *networkBus) GetBusListenerAndWriter() (chan busMessage, chan bool) {
	nB.currentCacheNum++
	return nB.busListeners[nB.currentCacheNum], nB.busWriters[nB.currentCacheNum]
}

func (nB *networkBus) AskModify(cacheNumber int, address Address) {
	bM :=busMessage{Address: address, LineState: Modify}
	nB.writeAndWaitOnBus(cacheNumber, bM)
}

func (nB *networkBus) AskExclusive(cacheNumber int, address Address) bool{
	bM :=busMessage{Address: address, LineState: Exclusive}

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

func (nB *networkBus) AskShared(cacheNumber int, address Address){
	bM :=busMessage{Address: address, LineState: Shared}
	nB.writeAndWaitOnBus(cacheNumber, bM)
}

func (nB *networkBus) writeAndWaitOnBus(cacheNumber int, bM busMessage) {
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