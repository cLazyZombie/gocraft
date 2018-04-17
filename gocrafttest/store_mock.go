package gocrafttest

import (
	. "github.com/cLazyZombie/gocraft/internal"
)

func NewStoreMock() *StoreMock {
	return &StoreMock{chunkBlocks: make(map[ChunkID][]BlockType)}
}

type StoreMock struct {
	chunkBlocks map[ChunkID][]BlockType
}

func (st *StoreMock) Add(bid BlockID, bt BlockType) {
	cid := bid.ChunkID()
	blocks, ok := st.chunkBlocks[cid]
	if !ok {
		blocks = make([]BlockType, ChunkWidth*ChunkWidth*ChunkWidth)
		st.chunkBlocks[cid] = blocks
	}

	blocks[bid.ToIndex()] = bt
}

func (st *StoreMock) ChunkBlocks(cid ChunkID) ([]BlockType, error) {
	var blocks []BlockType
	blocks, ok := st.chunkBlocks[cid]
	if !ok {
		return nil, nil
	}

	return blocks, nil
}

// func (st *StoreMock) RangeBlocks(id ChunkID, f func(bid BlockID, w BlockType)) error {
// 	bs, ok := st.chunkBlocks[id]
// 	if !ok {
// 		return nil
// 	}

// 	for _, b := range bs {
// 		f(b.bid, b.bt)
// 	}

// 	return nil
// }
