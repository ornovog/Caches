package Caches

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestMainMemory(t *testing.T) {
	var mM mainMemory
	mM.Init()

	mM.Store(0,5)
	x := mM.Fetch(0)
	assert.EqualValues(t,5, x)
}

func TestDirectMappedCache_GetData(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var dMC directMappedCache
	dMC.Init(&mM)
	address := uint32(0)
	collisionAddress := cacheSize + address

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(address,expectedVal)
	val, hit := dMC.GetData(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = dMC.GetData(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = dMC.GetData(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.EqualValues(t,false,hit)
	assert.EqualValues(t,expectedVal,mM.Fetch(address))

	val, hit = dMC.GetData(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.EqualValues(t,true,hit)

	val, hit = dMC.GetData(address)
	assert.EqualValues(t,val,expectedVal)
	assert.EqualValues(t,false,hit)
	assert.EqualValues(t,secondExpectedVal,mM.Fetch(collisionAddress))
}

func TestFullyAssociativeCache_GetData(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var fAC fullyAssociativeCache
	fAC.Init(&mM)
	address := uint32(0)
	collisionAddress := cacheSize + address

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(address,expectedVal)
	val, hit := fAC.GetData(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = fAC.GetData(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = fAC.GetData(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = fAC.GetData(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.EqualValues(t,true,hit)

	val, hit = fAC.GetData(address)
	assert.EqualValues(t,val,expectedVal)
	assert.EqualValues(t,true,hit)
}

func TestNWayAssociativeCache_GetData(t *testing.T) {
	var mM mainMemory
	mM.Init()

	var nWAC NWayAssociativeCache
	nWAC.Init(&mM)
	address := uint32(0)
	collisionAddress := cacheSize + address

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(address,expectedVal)
	val, hit := nWAC.GetData(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = nWAC.GetData(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(collisionAddress,secondExpectedVal)
	val, hit = nWAC.GetData(collisionAddress)
	assert.EqualValues(t,secondExpectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = nWAC.GetData(collisionAddress)
	assert.EqualValues(t,val,secondExpectedVal)
	assert.EqualValues(t,true,hit)

	val, hit = nWAC.GetData(address)
	assert.EqualValues(t,val,expectedVal)
	assert.EqualValues(t,true,hit)
}
