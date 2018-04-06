package main

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	ChunkWidth = 32
)

type BlockID struct {
	X, Y, Z int
}

func (v BlockID) Left() BlockID {
	return BlockID{v.X - 1, v.Y, v.Z}
}
func (v BlockID) Right() BlockID {
	return BlockID{v.X + 1, v.Y, v.Z}
}
func (v BlockID) Up() BlockID {
	return BlockID{v.X, v.Y + 1, v.Z}
}
func (v BlockID) Down() BlockID {
	return BlockID{v.X, v.Y - 1, v.Z}
}
func (v BlockID) Front() BlockID {
	return BlockID{v.X, v.Y, v.Z + 1}
}
func (v BlockID) Back() BlockID {
	return BlockID{v.X, v.Y, v.Z - 1}
}
func (v BlockID) ChunkID() ChunkID {
	return ChunkID{
		int(math.Floor(float64(v.X) / ChunkWidth)),
		0,
		int(math.Floor(float64(v.Z) / ChunkWidth)),
	}
}

func NearBlock(pos mgl32.Vec3) BlockID {
	return BlockID{
		int(round(pos.X())),
		int(round(pos.Y())),
		int(round(pos.Z())),
	}
}

// ChunkID : represent position of chunk
type ChunkID struct {
	X, Y, Z int
}

func (v ChunkID) Left() ChunkID {
	return ChunkID{v.X - 1, v.Y, v.Z}
}
func (v ChunkID) Right() ChunkID {
	return ChunkID{v.X + 1, v.Y, v.Z}
}
func (v ChunkID) Up() ChunkID {
	return ChunkID{v.X, v.Y + 1, v.Z}
}
func (v ChunkID) Down() ChunkID {
	return ChunkID{v.X, v.Y - 1, v.Z}
}
func (v ChunkID) Front() ChunkID {
	return ChunkID{v.X, v.Y, v.Z + 1}
}
func (v ChunkID) Back() ChunkID {
	return ChunkID{v.X, v.Y, v.Z - 1}
}

// Chunk : collection of blocks
type Chunk struct {
	id     ChunkID
	blocks sync.Map // map[Vec3]int

	Version int64
}

func NewChunk(id ChunkID) *Chunk {
	c := &Chunk{
		id:      id,
		Version: time.Now().Unix(),
	}
	return c
}

func (c *Chunk) UpdateVersion() {
	c.Version = time.Now().UnixNano() / int64(time.Millisecond)
}

func (c *Chunk) ID() ChunkID {
	return c.id
}

func (c *Chunk) Block(id BlockID) int {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}
	w, ok := c.blocks.Load(id)
	if ok {
		return w.(int)
	}
	return 0
}

func (c *Chunk) Add(id BlockID, w int) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}
	c.blocks.Store(id, w)
	c.UpdateVersion()
}

func (c *Chunk) Del(id BlockID) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}
	c.blocks.Delete(id)
	c.UpdateVersion()
}

func (c *Chunk) RangeBlocks(f func(id BlockID, w int)) {
	c.blocks.Range(func(key, value interface{}) bool {
		f(key.(BlockID), value.(int))
		return true
	})
}
