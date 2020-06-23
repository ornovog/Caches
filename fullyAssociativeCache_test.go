package caches

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestFullyAssociativeCache_Fetch(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var fAC fullyAssociativeCache
	fAC.Init(&mM)
	address := uint32(0)
	collisionAddress := cacheSize + address

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(address,expectedVal)
	val, hit := fAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.False(t,hit)

	val, hit = fAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = fAC.Fetch(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.False(t,hit)

	val, hit = fAC.Fetch(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.True(t,hit)

	val, hit = fAC.Fetch(address)
	assert.EqualValues(t,val,expectedVal)
	assert.True(t,hit)
}

func TestFullyAssociativeCache_Store(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var fAC fullyAssociativeCache
	fAC.Init(&mM)
	address := uint32(0)

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(address,expectedVal)
	hit := fAC.Store(address,expectedVal)
	assert.False(t,hit)

	val, hit := fAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)
}

func TestFullyAssociativeCache_Fetch_LRU(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var fAC fullyAssociativeCache
	fAC.Init(&mM)

	for line := 0; line < cacheSize ; line++ {
		fAC.Store(uint32(line),byte(line))
	}

	_, hit := fAC.Fetch(0)
	assert.True(t,hit)

	hit = fAC.Store(cacheSize,cacheSize)
	assert.False(t,hit)

	_, hit = fAC.Fetch(0)
	assert.True(t,hit)

	_, hit = fAC.Fetch(1)
	assert.False(t,hit)
}