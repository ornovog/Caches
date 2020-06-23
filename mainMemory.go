package caches

const memorySize = 4294967296 //2^32

type mainMemory struct {
	storage [memorySize]byte
}

func (mM *mainMemory) Init()  {
	mM.storage = [memorySize]byte{}
}

func (mM *mainMemory) Store(address uint32, data byte){
	mM.storage[address] = data
}

func (mM *mainMemory) Fetch(address uint32) byte{
	return mM.storage[address]
}