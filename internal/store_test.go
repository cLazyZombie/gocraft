package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_Encode_Decode_ChunkKey(t *testing.T) {
	cid := ChunkID{X: 1, Y: 2, Z: 3}
	b := encodeChunkDbKey(cid)
	assert.Equal(t, 12, len(b))
	assert.Equal(t, []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0}, b)

	decoded := decodeChunkDbKey(b)
	assert.Equal(t, cid, decoded)
}

func TestStore_Encode_Decode_ChunkValue(t *testing.T) {
	blocks := make([]BlockType, ChunkWidth*ChunkWidth*ChunkWidth)
	blocks[0] = 1
	blocks[1] = 2
	blocks[len(blocks)-1] = 9

	b, err := encodeChunkDbValue(blocks)
	assert.Nil(t, err)
	assert.Equal(t, ChunkWidth*ChunkWidth*ChunkWidth*2+1, len(b))
	assert.Equal(t, []byte{1, 1, 0, 2, 0}, b[:5])

	decoded := decodeChunkDbValue(b)
	assert.Equal(t, blocks, decoded)
}
