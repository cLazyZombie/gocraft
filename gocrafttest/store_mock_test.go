package gocrafttest

import (
	"testing"

	. "github.com/cLazyZombie/gocraft/internal"
	"github.com/stretchr/testify/assert"
)

func TestStoreMock(t *testing.T) {
	assert := assert.New(t)

	store := NewStoreMock()
	bidA := BlockID{X: 0, Y: 0, Z: ChunkWidth}
	store.Add(bidA, 1)

	callCount := 0
	store.RangeBlocks(bidA.ChunkID(), func(bid BlockID, bt BlockType) {
		callCount++
	})

	assert.Equal(1, callCount)

	callCount = 0
	store.RangeBlocks(ChunkID{X: 100, Y: 0, Z: 0}, func(bid BlockID, bt BlockType) {
		callCount++
	})

	assert.Equal(0, callCount)
}
