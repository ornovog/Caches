package caches

import "sync/atomic"

const memorySize = 1073741824 //2^32

type mainMemory struct {
	storage [memorySize]int32
}

func (mM *mainMemory) Init()  {
	mM.storage = [memorySize]int32{}
}

func (mM *mainMemory) Store(address uint32, data int32){
	atomic.StoreInt32(&mM.storage[address], data)
}

func (mM *mainMemory) Load(address uint32) int32{
	return atomic.LoadInt32(&mM.storage[address])
}