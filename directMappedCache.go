package Caches

const (
	indexBits = cacheSize - 1
	tagBits = addressMaxNumber - indexBits
)

type DMCacheLine struct {
	tag uint32
	data byte
}

type directMappedCache struct{
	storage [cacheSize]DMCacheLine
	mM *mainMemory
}

func (dMC *directMappedCache) GetData(address uint32) (byte, bool){
	index, tag, line := dMC.extractIndexTagAndLine(address)

	if line.tag == tag{
		return line.data, true
	}else {
		dMC.mM.Store(line.tag+index,line.data)
		data := dMC.mM.Fetch(address)
		line.data = data
		line.tag = tag

		return data, false
	}
}

func (dMC *directMappedCache) Update(address uint32, newData byte) bool{
	index, tag, line := dMC.extractIndexTagAndLine(address)

	if line.tag == tag {
		line.data = newData
		return true
	} else {
		dMC.mM.Store(line.tag+index, line.data)
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