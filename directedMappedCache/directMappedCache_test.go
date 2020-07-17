package directedMappedCache

import (
	"Caches"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestDirectMappedCache_Load(t *testing.T) {
	var mM caches.mainMemory
	mM.Init()

	var dMC DirectMappedCache
	dMC.Init(&mM)

	address := uint32(0)
	collisionAddress := caches.CacheSize + address

	expectedVal := int32(rand.Int())
	mM.Store(address,expectedVal)
	val, hit := dMC.Load(address)
	assert.EqualValues(t,expectedVal,val)
	assert.False(t,hit)

	val, hit = dMC.Load(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)

	secondExpectedVal := int32(rand.Int())
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = dMC.Load(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.False(t,hit)
	assert.EqualValues(t,expectedVal,mM.Load(address))

	val, hit = dMC.Load(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.True(t,hit)

	val, hit = dMC.Load(address)
	assert.EqualValues(t,val,expectedVal)
	assert.False(t,hit)
	assert.EqualValues(t,secondExpectedVal,mM.Load(collisionAddress))
}

func TestDirectMappedCache_Store(t *testing.T) {
	var mM caches.mainMemory
	mM.Init()

	var dMC DirectMappedCache
	dMC.Init(&mM)
	address := uint32(0)
	collisionAddress := caches.CacheSize + address

	expectedVal := int32(rand.Int())
	hit := dMC.Store(address,expectedVal)
	assert.False(t,hit)

	val, hit := dMC.Load(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)

	secondExpectedVal := int32(rand.Int())
	hit = dMC.Store(collisionAddress, secondExpectedVal)
	assert.False(t,hit)

	val, hit = dMC.Load(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.True(t,hit)

	val, hit = dMC.Load(address)
	assert.EqualValues(t,val,expectedVal)
	assert.False(t,hit)
}
