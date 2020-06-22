package Caches

import (
	"github.com/stretchr/testify/assert"
	"math"
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

	expectedVal := byte(rand.Intn(math.MaxInt8+1))
	mM.Store(address,expectedVal)
	val, hit := nWAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,false,hit)

	val, hit = nWAC.Fetch(address)
	assert.EqualValues(t,expectedVal,val)
	assert.EqualValues(t,true,hit)

	secondExpectedVal := byte(rand.Intn(math.MaxInt8+1))
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