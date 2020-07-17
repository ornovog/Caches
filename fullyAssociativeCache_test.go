package caches

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestFullyAssociativeCache_Load(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var fAC fullyAssociativeCache
	fAC.Init(&mM)
	address := Address(0)
	collisionAddress := CacheSize + address

	expectedVal := int32(rand.Int())
	mM.Store(address,expectedVal)
	val, hit := fAC.Load(address)
	assert.EqualValues(t,expectedVal,val)
	assert.False(t,hit)

	val, hit = fAC.Load(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)

	secondExpectedVal := int32(rand.Int())
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = fAC.Load(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.False(t,hit)

	val, hit = fAC.Load(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.True(t,hit)

	val, hit = fAC.Load(address)
	assert.EqualValues(t,val,expectedVal)
	assert.True(t,hit)
}

func TestFullyAssociativeCache_Store(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var fAC fullyAssociativeCache
	fAC.Init(&mM)
	address := uint32(0)

	expectedVal := int32(rand.Int())
	mM.Store(address,expectedVal)
	hit := fAC.Store(address,expectedVal)
	assert.False(t,hit)

	val, hit := fAC.Load(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)
}

func TestFullyAssociativeCache_Load_LRU(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var fAC fullyAssociativeCache
	fAC.Init(&mM)

	for line := 0; line < CacheSize; line++ {
		fAC.Store(uint32(line),int32(line))
	}

	_, hit := fAC.Load(0)
	assert.True(t,hit)

	hit = fAC.Store(CacheSize, CacheSize)
	assert.False(t,hit)

	_, hit = fAC.Load(0)
	assert.True(t,hit)

	_, hit = fAC.Load(1)
	assert.False(t,hit)
}