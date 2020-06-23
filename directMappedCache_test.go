package caches

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestDirectMappedCache_Fetch(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var dMC directMappedCache
	dMC.Init(&mM)

	address := uint32(0)
	collisionAddress := cacheSize + address

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(address,expectedVal)
	val, hit := dMC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.False(t,hit)

	val, hit = dMC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = dMC.Fetch(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.False(t,hit)
	assert.EqualValues(t,expectedVal,mM.Fetch(address))

	val, hit = dMC.Fetch(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.True(t,hit)

	val, hit = dMC.Fetch(address)
	assert.EqualValues(t,val,expectedVal)
	assert.False(t,hit)
	assert.EqualValues(t,secondExpectedVal,mM.Fetch(collisionAddress))
}

func TestDirectMappedCache_Store(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var dMC directMappedCache
	dMC.Init(&mM)
	address := uint32(0)
	collisionAddress := cacheSize + address

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	hit := dMC.Store(address,expectedVal)
	assert.False(t,hit)

	val, hit := dMC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.True(t,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	hit = dMC.Store(collisionAddress, secondExpectedVal)
	assert.False(t,hit)

	val, hit = dMC.Fetch(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.True(t,hit)

	val, hit = dMC.Fetch(address)
	assert.EqualValues(t,val,expectedVal)
	assert.False(t,hit)
}
