package Caches

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
	assert.EqualValues(t,false,hit)

	val, hit = dMC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = dMC.Fetch(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.EqualValues(t,false,hit)
	assert.EqualValues(t,expectedVal,mM.Fetch(address))

	val, hit = dMC.Fetch(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.EqualValues(t,true,hit)

	val, hit = dMC.Fetch(address)
	assert.EqualValues(t,val,expectedVal)
	assert.EqualValues(t,false,hit)
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
	mM.Store(address,expectedVal)
	hit := dMC.Store(address,expectedVal)
	assert.EqualValues(t,false,hit)

	val, hit := dMC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	hit = dMC.Store(collisionAddress, secondExpectedVal)
	assert.EqualValues(t,false,hit)

	val, hit = dMC.Fetch(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.EqualValues(t,true,hit)

	val, hit = dMC.Fetch(address)
	assert.EqualValues(t,val,expectedVal)
	assert.EqualValues(t,false,hit)
}
