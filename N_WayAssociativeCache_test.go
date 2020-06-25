package caches

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNWayAssociativeCache_Fetch(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var nWAC NWayAssociativeCache
	nWAC.Init(&mM)
	address := uint32(0)
	collisionAddress := cacheSize + address

	expectedVal := int32(rand.Int())
	mM.Store(address,expectedVal)
	val, hit := nWAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = nWAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)

	secondExpectedVal := int32(rand.Int())
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = nWAC.Fetch(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = nWAC.Fetch(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.EqualValues(t,true,hit)

	val, hit = nWAC.Fetch(address)
	assert.EqualValues(t,val,expectedVal)
	assert.EqualValues(t,true,hit)
}

func TestNWayAssociativeCache_Store(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var nWAC NWayAssociativeCache
	nWAC.Init(&mM)
	address := uint32(0)

	expectedVal := int32(rand.Int())
	mM.Store(address,expectedVal)
	hit := nWAC.Store(address,expectedVal)
	assert.EqualValues(t,false,hit)

	val, hit := nWAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)
}

func TestNWayAssociativeCache_Fetch_LRU(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var nWAC NWayAssociativeCache
	nWAC.Init(&mM)

	for line := 0; line < cacheSize ; line+=numOfWays {
		nWAC.Store(uint32(line),int32(line))
	}

	_, hit := nWAC.Fetch(0)
	assert.True(t,hit)

	hit = nWAC.Store(cacheSize,cacheSize)
	assert.False(t,hit)

	_, hit = nWAC.Fetch(0)
	assert.True(t,hit)

	_, hit = nWAC.Fetch(numOfWays)
	assert.False(t,hit)
}