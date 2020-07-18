package mainMemory

const memorySize = 1073741824 //2^32

type Address uint32
type Data int32

type MainMemory struct {
	storage [memorySize]Data
}

func (mM *MainMemory) Init() {
	mM.storage = [memorySize]Data{}
}

func (mM *MainMemory) Store(address Address, data Data) {
	mM.storage[address] = data
}

func (mM *MainMemory) Load(address Address) Data {
	return mM.storage[address]
}
