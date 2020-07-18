package fullyAssociativeCache

import (
	caches "Caches"
	"Caches/bus"
	"Caches/mainMemory"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestFullyAssociativeCache_Load(t *testing.T) {
	var mM mainMemory.MainMemory
	mM.Init()

	var nB bus.NetworkBus
	nB.Init(1)

	var fAC fullyAssociativeCache
	fAC.Init(&mM, &nB)
	address := mainMemory.Address(0)
	collisionAddress := caches.CacheSize + address

	expectedVal := mainMemory.Data(rand.Int())
	mM.Store(address, expectedVal)
	val, hit := fAC.Load(address)
	assert.EqualValues(t, expectedVal, val)
	assert.False(t, hit)

	val, hit = fAC.Load(address)
	assert.EqualValues(t, expectedVal, val)
	assert.True(t, hit)

	secondExpectedVal := mainMemory.Data(rand.Int())
	mM.Store(collisionAddress, secondExpectedVal)
	val, hit = fAC.Load(collisionAddress)
	assert.EqualValues(t, secondExpectedVal, val)
	assert.False(t, hit)

	val, hit = fAC.Load(collisionAddress)
	assert.EqualValues(t, val, secondExpectedVal)
	assert.True(t, hit)

	val, hit = fAC.Load(address)
	assert.EqualValues(t, val, expectedVal)
	assert.True(t, hit)
}

func TestFullyAssociativeCache_Store(t *testing.T) {
	var mM mainMemory.MainMemory
	mM.Init()

	var nB bus.NetworkBus
	nB.Init(1)

	var fAC fullyAssociativeCache
	fAC.Init(&mM, &nB)
	address := mainMemory.Address(0)

	expectedVal := mainMemory.Data(rand.Int())
	mM.Store(address, expectedVal)
	hit := fAC.Store(address, expectedVal)
	assert.False(t, hit)

	val, hit := fAC.Load(address)
	assert.EqualValues(t, expectedVal, val)
	assert.True(t, hit)
}

func TestFullyAssociativeCache_Load_LRU(t *testing.T) {
	var mM mainMemory.MainMemory
	mM.Init()

	var nB bus.NetworkBus
	nB.Init(1)

	var fAC fullyAssociativeCache
	fAC.Init(&mM, &nB)

	for line := 0; line < caches.CacheSize; line++ {
		fAC.Store(mainMemory.Address(line), mainMemory.Data(line))
	}

	_, hit := fAC.Load(0)
	assert.True(t, hit)

	hit = fAC.Store(caches.CacheSize, caches.CacheSize)
	assert.False(t, hit)

	_, hit = fAC.Load(0)
	assert.True(t, hit)

	_, hit = fAC.Load(1)
	assert.False(t, hit)
}
