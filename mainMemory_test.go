package caches

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMainMemory(t *testing.T) {
	var mM mainMemory
	mM.Init()

	mM.Store(0,5)
	x := mM.Load(0)
	assert.EqualValues(t,5, x)
}





