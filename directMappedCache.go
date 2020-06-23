package caches

const (
	indexBits = cacheSize - 1
	tagBits = addressMaxNumber - indexBits
)

//DMCacheLine - Direct Mapped Cache Line
type DMCacheLine struct {
	valid bool
	tag uint32
	data byte
}

type directMappedCache struct{
	storage [cacheSize]DMCacheLine
	mainMemory *mainMemory
}

func (dMC *directMappedCache)Init(mainMemory *mainMemory){
	dMC.storage = [cacheSize]DMCacheLine{}
	dMC.mainMemory = mainMemory
}

func (dMC *directMappedCache) Fetch(address uint32) (byte, bool){
	index, tag:= dMC.extractIndexAndTag(address)
	line :=dMC.storage[index]

	if  line.valid{
		if line.tag == tag {
			return line.data, true
		}
		dMC.mainMemory.Store(line.tag+index,line.data)
	}

	data := dMC.mainMemory.Fetch(address)

	dMC.storage[index].data = data
	dMC.storage[index].tag = tag
	dMC.storage[index].valid = true

	return data, false
}

func (dMC *directMappedCache) Store(address uint32, newData byte) bool{
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

func (dMC *directMappedCache) extractIndexAndTag(address uint32) (uint32, uint32) {
	index := address & indexBits
	tag := address & tagBits
	return index, tag
}