package main

import (
	"log"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	lru "github.com/hashicorp/golang-lru"
)

type World struct {
	mutex  sync.Mutex
	chunks *lru.Cache // map[ChunkID]*Chunk
	store  IStore
}

func NewWorld(store IStore) *World {
	m := (*renderRadius) * (*renderRadius) * (*renderRadius) * 4
	chunks, _ := lru.New(m)
	return &World{
		chunks: chunks,
		store:  store,
	}
}

func (w *World) loadChunk(id ChunkID) (*Chunk, bool) {
	chunk, ok := w.chunks.Get(id)
	if !ok {
		return nil, false
	}
	return chunk.(*Chunk), true
}

func (w *World) storeChunk(id ChunkID, chunk *Chunk) {
	w.chunks.Add(id, chunk)
}

func (w *World) Collide(pos mgl32.Vec3) (mgl32.Vec3, bool) {
	x, y, z := pos.X(), pos.Y(), pos.Z()
	nx, ny, nz := round(pos.X()), round(pos.Y()), round(pos.Z())
	const pad = 0.25

	head := BlockID{int(nx), int(ny), int(nz)}
	foot := head.Down()

	stop := false
	for _, b := range []BlockID{foot, head} {
		if w.Block(b.Left()).IsObstacle() && x < nx && nx-x > pad {
			x = nx - pad
		}
		if w.Block(b.Right()).IsObstacle() && x > nx && x-nx > pad {
			x = nx + pad
		}
		if w.Block(b.Down()).IsObstacle() && y < ny && ny-y > pad {
			y = ny - pad
			stop = true
		}
		if w.Block(b.Up()).IsObstacle() && y > ny && y-ny > pad {
			y = ny + pad
			stop = true
		}
		if w.Block(b.Back()).IsObstacle() && z < nz && nz-z > pad {
			z = nz - pad
		}
		if w.Block(b.Front()).IsObstacle() && z > nz && z-nz > pad {
			z = nz + pad
		}
	}
	return mgl32.Vec3{x, y, z}, stop
}

func (w *World) HitTest(pos mgl32.Vec3, vec mgl32.Vec3) (*BlockID, *BlockID) {
	var (
		maxLen = float32(8.0)
		step   = float32(0.125)

		block, prev BlockID
		pprev       *BlockID
	)

	for len := float32(0); len < maxLen; len += step {
		block = NearBlock(pos.Add(vec.Mul(len)))
		if prev != block && w.HasBlock(block) {
			return &block, pprev
		}
		prev = block
		pprev = &prev
	}
	return nil, nil
}

func (w *World) Block(id BlockID) BlockType {
	chunk := w.BlockChunk(id)
	if chunk == nil {
		return 0
	}
	return chunk.Block(id)
}

func (w *World) BlockChunk(bid BlockID) *Chunk {
	cid := bid.ChunkID()
	chunk, ok := w.loadChunk(cid)
	if !ok {
		return nil
	}
	return chunk
}

func (w *World) HasBlock(bid BlockID) bool {
	tp := w.Block(bid)
	return tp != 0
}

func (w *World) Chunk(cid ChunkID) *Chunk {
	p, ok := w.loadChunk(cid)
	if ok {
		return p
	}
	chunk := NewChunk(cid)
	blocks := makeChunkMap(cid)
	for block, tp := range blocks {
		chunk.Add(block, tp)
	}
	err := w.store.RangeBlocks(cid, func(bid BlockID, w BlockType) {
		if w == 0 {
			chunk.Del(bid)
			return
		}
		chunk.Add(bid, w)
	})
	if err != nil {
		log.Printf("fetch chunk(%v) from db error:%s", cid, err)
		return nil
	}
	w.storeChunk(cid, chunk)
	return chunk
}

func (w *World) Chunks(cids []ChunkID) []*Chunk {
	ch := make(chan *Chunk)
	var chunks []*Chunk
	for _, cid := range cids {
		go func(cid ChunkID) {
			ch <- w.Chunk(cid)
		}(cid)
	}
	for range cids {
		chunk := <-ch
		if chunk != nil {
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func makeChunkMap(cid ChunkID) map[BlockID]BlockType {
	const (
		grassBlock = 1
		sandBlock  = 2
		grass      = 17
		leaves     = 15
		wood       = 5
	)
	m := make(map[BlockID]BlockType)
	startY, endY := cid.Y*ChunkWidth, (cid.Y+1)*ChunkWidth-1
	p, q := cid.X, cid.Z
	for dx := 0; dx < ChunkWidth; dx++ {
		for dz := 0; dz < ChunkWidth; dz++ {
			x, z := p*ChunkWidth+dx, q*ChunkWidth+dz
			f := noise2(float32(x)*0.01, float32(z)*0.01, 4, 0.5, 2)
			g := noise2(float32(-x)*0.01, float32(-z)*0.01, 2, 0.9, 2)
			mh := int(g*32 + 16)
			h := int(f * float32(mh))
			var w BlockType = grassBlock
			if h <= 12 {
				h = 12
				w = sandBlock
			}

			// grass and sand
			for y := 0; y < h; y++ {
				if y >= startY && y <= endY {
					m[BlockID{x, y, z}] = w
				}
			}

			// flowers
			if h >= startY && h <= endY {
				if w == grassBlock {
					if noise2(-float32(x)*0.1, float32(z)*0.1, 4, 0.8, 2) > 0.6 {
						m[BlockID{x, h, z}] = grass
					}
					if noise2(float32(x)*0.05, float32(-z)*0.05, 4, 0.8, 2) > 0.7 {
						w := BlockType(18 + int(noise2(float32(x)*0.1, float32(z)*0.1, 4, 0.8, 2)*7))
						m[BlockID{x, h, z}] = w
					}
				}
			}

			// tree
			if w == 1 {
				ok := true
				if dx-4 < 0 || dz-4 < 0 ||
					dx+4 > ChunkWidth || dz+4 > ChunkWidth {
					ok = false
				}
				if ok && noise2(float32(x), float32(z), 6, 0.5, 2) > 0.79 {
					for y := h + 3; y < h+8; y++ {
						for ox := -3; ox <= 3; ox++ {
							for oz := -3; oz <= 3; oz++ {
								d := ox*ox + oz*oz + (y-h-4)*(y-h-4)
								if d < 11 {
									if y >= startY && y <= endY {
										m[BlockID{x + ox, y, z + oz}] = leaves
									}
								}
							}
						}
					}
					for y := h; y < h+7; y++ {
						if y >= startY && y <= endY {
							m[BlockID{x, y, z}] = wood
						}
					}
				}
			}

			// cloud
			for y := 64; y < 72; y++ {
				if y >= startY && y <= endY && noise3(float32(x)*0.01, float32(y)*0.1, float32(z)*0.01, 8, 0.5, 2) > 0.69 {
					m[BlockID{x, y, z}] = 16
				}
			}
		}
	}
	return m
}
