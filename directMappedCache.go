package Caches

const (
	indexBits = cacheSize - 1
	tagBits = addressMaxNumber - indexBits
)

type DMCacheLine struct {
	valid bool
	tag uint32
	data byte
}

type directMappedCache struct{
	storage *[cacheSize]DMCacheLine
	mainMemory *mainMemory
}

func (dMC *directMappedCache)Init(mainMemory *mainMemory){
	dMC.storage = &[cacheSize]DMCacheLine{}
	dMC.mainMemory = mainMemory
}

func (dMC *directMappedCache) GetData(address uint32) (byte, bool){
	index, tag, line := dMC.extractIndexTagAndLine(address)

	if line.tag == tag && line.valid{
		return line.data, true
	}else {
		if line.valid{
			dMC.mainMemory.Store(line.tag+index,line.data)
		}

		data := dMC.mainMemory.Fetch(address)
		line.data = data
		line.tag = tag
		line.valid = true

		return data, false
	}
}

func (dMC *directMappedCache) Update(address uint32, newData byte) bool{
	index, tag, line := dMC.extractIndexTagAndLine(address)

	if line.tag == tag {
		line.data = newData
		return true
	} else {
		dMC.mainMemory.Store(line.tag+index, line.data)
		line.data = newData
		line.tag = tag

		return false
	}
}

func (dMC *directMappedCache) extractIndexTagAndLine(address uint32) (uint32, uint32, *DMCacheLine) {
	index := address & indexBits
	tag := address & tagBits
	line := &dMC.storage[index]
	return index, tag, line
}