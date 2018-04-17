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

	blocks, err := store.ChunkBlocks(bidA.ChunkID())
	assert.Nil(err)

	assert.Equal(BlockType(1), blocks[bidA.ToIndex()])
}
