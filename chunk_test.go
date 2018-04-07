package main

import (
	"testing"
)

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
	vals := make(map[BlockID]int)

	chunk := NewChunk(ChunkID{1, 0, 1})
	f := func(id BlockID, w int) {
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
		t.Fatalf("f should be called 2 times")
	}
}
