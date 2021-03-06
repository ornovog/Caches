package N_WayAssociativeCache

import (
	"Caches"
	"Caches/mainMemory"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNWayAssociativeCache_Fetch(t *testing.T) {
	var mM mainMemory.MainMemory
	mM.Init()

	var nWAC NWayAssociativeCache
	nWAC.Init(&mM)
	address := uint32(0)
	collisionAddress := caches.CacheSize + address

	expectedVal := int32(rand.Int())
	mM.Store(address, expectedVal)
	val, hit := nWAC.Fetch(address)
	assert.EqualValues(t, expectedVal, val)
	assert.EqualValues(t, false, hit)

	val, hit = nWAC.Fetch(address)
	assert.EqualValues(t, expectedVal, val)
	assert.EqualValues(t, true, hit)

	secondExpectedVal := int32(rand.Int())
	mM.Store(collisionAddress, secondExpectedVal)
	val, hit = nWAC.Fetch(collisionAddress)
	assert.EqualValues(t, secondExpectedVal, val)
	assert.EqualValues(t, false, hit)

	val, hit = nWAC.Fetch(collisionAddress)
	assert.EqualValues(t, val, secondExpectedVal)
	assert.EqualValues(t, true, hit)

	val, hit = nWAC.Fetch(address)
	assert.EqualValues(t, val, expectedVal)
	assert.EqualValues(t, true, hit)
}

func TestNWayAssociativeCache_Store(t *testing.T) {
	var mM mainMemory.MainMemory
	mM.Init()

	var nWAC NWayAssociativeCache
	nWAC.Init(&mM)
	address := uint32(0)

	expectedVal := int32(rand.Int())
	mM.Store(address, expectedVal)
	hit := nWAC.Store(address, expectedVal)
	assert.EqualValues(t, false, hit)

	val, hit := nWAC.Fetch(address)
	assert.EqualValues(t, expectedVal, val)
	assert.EqualValues(t, true, hit)
}

func TestNWayAssociativeCache_Fetch_LRU(t *testing.T) {
	var mM mainMemory.MainMemory
	mM.Init()

	var nWAC NWayAssociativeCache
	nWAC.Init(&mM)

	for line := 0; line < caches.CacheSize; line += NumOfWays {
		nWAC.Store(uint32(line), int32(line))
	}

	_, hit := nWAC.Fetch(0)
	assert.True(t, hit)

	hit = nWAC.Store(caches.CacheSize, caches.CacheSize)
	assert.False(t, hit)

	_, hit = nWAC.Fetch(0)
	assert.True(t, hit)

	_, hit = nWAC.Fetch(NumOfWays)
	assert.False(t, hit)
}
