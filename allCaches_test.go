package Caches

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func initializeMainMemory() *mainMemory{
	mM := mainMemory{storage: make([]byte,memorySize)}

	for i,_ := range mM.storage{
		mM.storage[i] = byte(rand.Intn(math.MaxInt8+1))
	}
	return &mM
}

func TestMainMemory(t *testing.T) {
	var mM mainMemory
	mM.Store(0,5)
	x := mM.Fetch(0)
	assert.EqualValues(t,5, x)
}

func TestDirectMappedCache_GetData(t *testing.T) {
	mM := initializeMainMemory()
	dMC := directMappedCache{mM: mM}

	v, hit := dMC.GetData(0)
	assert.EqualValues(t,v,mM.storage[0])
	assert.EqualValues(t,hit,false)

	v, hit = dMC.GetData(0)
	assert.EqualValues(t,v,mM.storage[0])
	assert.EqualValues(t,hit,false)
}

func TestFullyAssociativeCache_GetData(t *testing.T) {

}

func TestNWayAssociativeCache_GetData(t *testing.T) {

}
