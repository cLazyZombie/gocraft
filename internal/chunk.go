package internal

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
		int(math.Floor(float64(v.Y) / ChunkWidth)),
		int(math.Floor(float64(v.Z) / ChunkWidth)),
	}
}
func (v BlockID) ToIndex() int {
	if v.X < 0 {
		v.X = ChunkWidth + v.X%ChunkWidth
	}

	if v.Y < 0 {
		v.Y = ChunkWidth + v.Y%ChunkWidth
	}

	if v.Z < 0 {
		v.Z = ChunkWidth + v.Z%ChunkWidth
	}

	return (v.X % ChunkWidth) + (v.Y%ChunkWidth)*ChunkWidth + (v.Z%ChunkWidth)*ChunkWidth*ChunkWidth
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
	id ChunkID

	blocks []BlockType
	locker sync.Locker

	Version int64
}

func NewChunk(id ChunkID) *Chunk {
	c := &Chunk{
		id:      id,
		locker:  &sync.Mutex{},
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

func (c *Chunk) Block(id BlockID) BlockType {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	if len(c.blocks) == 0 {
		return 0
	}

	return c.blocks[id.ToIndex()]
}

func (c *Chunk) Add(id BlockID, w BlockType) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	// add empty ceblock into empty chunk do nothing
	if len(c.blocks) == 0 && w == 0 {
		return
	}

	if len(c.blocks) == 0 {
		c.blocks = make([]BlockType, ChunkWidth*ChunkWidth*ChunkWidth)
	}

	c.blocks[id.ToIndex()] = w
	c.UpdateVersion()
}

func (c *Chunk) Del(id BlockID) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	if len(c.blocks) == 0 {
		log.Panicln("Del to empth block")
		return
	}

	c.blocks[id.ToIndex()] = 0

	// todo: delete blocks if there is no visible block

	c.UpdateVersion()
}

func (c *Chunk) RangeBlocks(f func(id BlockID, w BlockType)) {
	if len(c.blocks) == 0 {
		return
	}

	sx, sy, sz := c.id.X*ChunkWidth, c.id.Y*ChunkWidth, c.id.Z*ChunkWidth
	for z := 0; z < ChunkWidth; z++ {
		for y := 0; y < ChunkWidth; y++ {
			for x := 0; x < ChunkWidth; x++ {
				id := BlockID{x + sx, y + sy, z + sz}

				c.locker.Lock()
				w := c.blocks[id.ToIndex()]
				c.locker.Unlock()

				if w != 0 {
					f(id, w)
				}
			}
		}
	}
}
