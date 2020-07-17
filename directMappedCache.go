package caches

const (
	indexBits = CacheSize - 1
	tagBits = AddressMaxNumber - indexBits
)

//DMCacheLine - Direct Mapped Cache Line
type DMCacheLine struct {
	valid bool
	tag Address
	data Data
}

type directMappedCache struct{
	storage [CacheSize]DMCacheLine
	mainMemory *mainMemory
}

func (dMC *directMappedCache)Init(mainMemory *mainMemory){
	dMC.storage = [CacheSize]DMCacheLine{}
	dMC.mainMemory = mainMemory
}

func (dMC *directMappedCache) Load(address Address) (Data, bool){
	index, tag:= dMC.extractIndexAndTag(address)
	line :=dMC.storage[index]

	if  line.valid{
		if line.tag == tag {
			return line.data, true
		}
		dMC.mainMemory.Store(line.tag+index,line.data)
	}

	data := dMC.mainMemory.Load(address)

	dMC.storage[index].data = data
	dMC.storage[index].tag = tag
	dMC.storage[index].valid = true

	return data, false
}

func (dMC *directMappedCache) Store(address Address, newData Data) bool{
	index, tag := dMC.extractIndexAndTag(address)
	line := dMC.storage[index]

	if line.valid {
		if line.tag == tag {
			dMC.storage[index].data = newData
			return true
		}
		dMC.mainMemory.Store(line.tag+index, line.data)
	}

	dMC.storage[index].data = newData
	dMC.storage[index].tag = tag
	dMC.storage[index].valid = true

	return false
}

func (dMC *directMappedCache) extractIndexAndTag(address Address) (Address, Address) {
	index := address & indexBits
	tag := address & tagBits
	return index, tag
}