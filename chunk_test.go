package main

import (
	"testing"
)

func TestBlockID_ToIndex(t *testing.T) {
	id := BlockID{ChunkWidth + 1, ChunkWidth + 1, ChunkWidth + 1}

	if id.ToIndex() != 1+ChunkWidth+ChunkWidth*ChunkWidth {
		t.Fatalf("expected %d, but actual: %d", 1+ChunkWidth+ChunkWidth*ChunkWidth, id.ToIndex())
	}

	id = BlockID{-1, 0, 0}
	if id.ToIndex() != ChunkWidth-1 {
		t.Fatalf("expected %d, but actual: %d", ChunkWidth-1, id.ToIndex())
	}

	id = BlockID{-32, 0, 0}
	if id.ToIndex() != 0 {
		t.Fatalf("expected %d, but actual: %d", 0, id.ToIndex())
	}

	id = BlockID{-33, 0, 0}
	if id.ToIndex() != ChunkWidth-1 {
		t.Fatalf("expected %d, but actual: %d", ChunkWidth-1, id.ToIndex())
	}

	id = BlockID{-64, 0, 0}
	if id.ToIndex() != 0 {
		t.Fatalf("expected %d, but actual: %d", 0, id.ToIndex())
	}
}

func TestChunk_Add(t *testing.T) {
	chunk := NewChunk(ChunkID{1, 0, 1})

	blockID := BlockID{ChunkWidth + 1, 0, ChunkWidth + 1}
	block := chunk.Block(blockID)
	if block != 0 {
		t.Fatalf("block should be empty")
	}

	chunk.Add(blockID, 1)
	block = chunk.Block(blockID)
	if block == 0 {
		t.Fatalf("block should be filled")
	}
}

func TestChunk_Del(t *testing.T) {
	chunk := NewChunk(ChunkID{1, 0, 1})
	blockID := BlockID{ChunkWidth + 1, 0, ChunkWidth + 1}
	chunk.Add(blockID, 1)
	chunk.Del(blockID)

	block := chunk.Block(blockID)
	if block != 0 {
		t.Fatalf("block should be removed")
	}
}

func TestChunk_RangeBlocks(t *testing.T) {
	vals := make(map[BlockID]BlockType)

	chunk := NewChunk(ChunkID{1, 0, 1})
	f := func(id BlockID, w BlockType) {
		vals[id] = w
	}

	chunk.RangeBlocks(f)

	if len(vals) != 0 {
		t.Fatalf("Range should not be called")
	}

	blockA := BlockID{ChunkWidth + 1, 0, ChunkWidth + 1}
	blockB := BlockID{ChunkWidth + 2, 0, ChunkWidth + 2}
	chunk.Add(blockA, 1)
	chunk.Add(blockB, 2)

	chunk.RangeBlocks(f)
	if len(vals) != 2 {
		t.Fatalf("f should be called 2 times. actual: %d", len(vals))
	}
}
