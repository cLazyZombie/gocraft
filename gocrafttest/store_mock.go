package gocrafttest

import . "github.com/cLazyZombie/gocraft/internal"

func NewStoreMock() *StoreMock {
	return &StoreMock{chunkBlocks: make(map[ChunkID][]StoredBlock)}
}

type StoreMock struct {
	chunkBlocks map[ChunkID][]StoredBlock
}

type StoredBlock struct {
	bid BlockID
	bt  BlockType
}

func (st *StoreMock) Add(bid BlockID, bt BlockType) {
	cid := bid.ChunkID()
	blocks, ok := st.chunkBlocks[cid]
	if !ok {
		st.chunkBlocks[cid] = make([]StoredBlock, 0)
	}

	blocks = append(blocks, StoredBlock{bid: bid, bt: bt})
	st.chunkBlocks[cid] = blocks
}

func (st *StoreMock) RangeBlocks(id ChunkID, f func(bid BlockID, w BlockType)) error {
	bs, ok := st.chunkBlocks[id]
	if !ok {
		return nil
	}

	for _, b := range bs {
		f(b.bid, b.bt)
	}

	return nil
}
