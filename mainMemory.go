package caches

const memorySize = 1073741824 //2^32

type Address uint32
type Data int32

type mainMemory struct {
	storage [memorySize]Data
}

func (mM *mainMemory) Init()  {
	mM.storage = [memorySize]Data{}
}

func (mM *mainMemory) Store(address Address, data Data){
	mM.storage[address] = data
}

func (mM *mainMemory) Load(address Address) Data{
	return mM.storage[address]
}