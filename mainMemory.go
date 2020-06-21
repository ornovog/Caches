package Caches

const memorySize = 4294967296 //2^32

type mainMemory struct {
	storage [memorySize]float64
}

func (mM *mainMemory) Store(address uint32, data float64){
	mM.storage[address] = data
}

func (mM *mainMemory) Fetch(address uint32) float64{
	return  mM.storage[address]
}