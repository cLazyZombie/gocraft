package internal

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"

	"github.com/boltdb/bolt"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	dbpath = flag.String("db", "gocraft.db", "db file name")
)

var (
	//blockBucket  = []byte("block")
	chunkBucket  = []byte("chunk")
	cameraBucket = []byte("camera")

	GlobalStore *Store
)

func InitStore() error {
	if *dbpath == "" {
		return nil
	}
	var err error
	GlobalStore, err = NewStore(*dbpath)
	return err
}

type IStore interface {
	ChunkBlocks(cid ChunkID) ([]BlockType, error)
}

type Store struct {
	db *bolt.DB
}

func NewStore(p string) (*Store, error) {
	db, err := bolt.Open(p, 0666, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(chunkBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(cameraBucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	db.NoSync = true
	return &Store{
		db: db,
	}, nil
}

func (s *Store) UpdateChunk(cid ChunkID, blocks []BlockType) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		log.Printf("put chunk[%d]", cid)
		bkt := tx.Bucket(chunkBucket)
		key := encodeChunkDbKey(cid)
		value, err := encodeChunkDbValue(blocks)
		if err != nil {
			return err
		}
		return bkt.Put(key, value)
	})
}

// func (s *Store) UpdateBlock(bid BlockID, w BlockType) error {
// 	return s.db.Update(func(tx *bolt.Tx) error {
// 		log.Printf("put %v -> %d", bid, w)
// 		bkt := tx.Bucket(blockBucket)
// 		cid := bid.ChunkID()
// 		key := encodeBlockDbKey(cid, bid)
// 		value := encodeBlockDbValue(w)
// 		return bkt.Put(key, value)
// 	})
// }

func (s *Store) UpdateCamera(pos mgl32.Vec3, rx, ry float32) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(cameraBucket)
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, &pos)
		binary.Write(buf, binary.LittleEndian, [...]float32{rx, ry})
		bkt.Put(cameraBucket, buf.Bytes())
		return nil
	})
}

func (s *Store) GetCamera() (mgl32.Vec3, float32, float32) {
	var pos = mgl32.Vec3{0, 16, 0}
	var rx, ry float32
	s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(cameraBucket)
		value := bkt.Get(cameraBucket)
		if value == nil {
			return nil
		}
		buf := bytes.NewBuffer(value)
		binary.Read(buf, binary.LittleEndian, &pos)
		binary.Read(buf, binary.LittleEndian, &rx)
		binary.Read(buf, binary.LittleEndian, &ry)
		return nil
	})
	return pos, rx, ry
}

func (s *Store) ChunkBlocks(cid ChunkID) ([]BlockType, error) {
	var blocks []BlockType
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(chunkBucket)
		key := encodeChunkDbKey(cid)
		value := bkt.Get(key)
		if value == nil {
			return nil
		}

		blocks = decodeChunkDbValue(value)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return blocks, nil
}

// func (s *Store) RangeBlocks(id ChunkID, f func(bid BlockID, w BlockType)) error {
// 	return s.db.View(func(tx *bolt.Tx) error {
// 		bkt := tx.Bucket(blockBucket)
// 		startkey := encodeBlockDbKey(id, BlockID{0, 0, 0})
// 		iter := bkt.Cursor()
// 		for k, v := iter.Seek(startkey); k != nil; k, v = iter.Next() {
// 			cid, bid := decodeBlockDbKey(k)
// 			if cid != id {
// 				break
// 			}
// 			w := decodeBlockDbValue(v)
// 			f(bid, w)
// 		}
// 		return nil
// 	})
// }

func (s *Store) Close() {
	s.db.Sync()
	s.db.Close()
}

func encodeChunkDbKey(cid ChunkID) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, [...]int32{int32(cid.X), int32(cid.Y), int32(cid.Z)})
	return buf.Bytes()
}

// func encodeBlockDbKey(cid ChunkID, bid BlockID) []byte {
// 	buf := new(bytes.Buffer)
// 	binary.Write(buf, binary.LittleEndian, [...]int32{int32(cid.X), int32(cid.Z)})
// 	binary.Write(buf, binary.LittleEndian, [...]int32{int32(bid.X), int32(bid.Y), int32(bid.Z)})
// 	return buf.Bytes()
// }

func decodeChunkDbKey(b []byte) ChunkID {
	if len(b) != 4*3 {
		log.Panicf("abd db key length:%d", len(b))
	}
	buf := bytes.NewBuffer(b)
	var arr [3]int32
	binary.Read(buf, binary.LittleEndian, &arr)

	cid := ChunkID{int(arr[0]), int(arr[1]), int(arr[2])}
	return cid
}

// func decodeBlockDbKey(b []byte) (ChunkID, BlockID) {
// 	if len(b) != 4*5 {
// 		log.Panicf("bad db key length:%d", len(b))
// 	}
// 	buf := bytes.NewBuffer(b)
// 	var arr [5]int32
// 	binary.Read(buf, binary.LittleEndian, &arr)

// 	cid := ChunkID{int(arr[0]), 0, int(arr[1])}
// 	bid := BlockID{int(arr[2]), int(arr[3]), int(arr[4])}
// 	if bid.ChunkID() != cid {
// 		log.Panicf("bad db key: cid:%v, bid:%v", cid, bid)
// 	}
// 	return cid, bid
// }

func encodeChunkDbValue(blocks []BlockType) ([]byte, error) {
	// first byte of value indicate emptyness of chunk (0 == empty, other = has block)

	if len(blocks) == 0 {
		return []byte{0}, nil
	}

	value := make([]byte, 0, (ChunkWidth*ChunkWidth*ChunkWidth)*2+1)
	buff := bytes.NewBuffer(value)
	err := buff.WriteByte(1)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.LittleEndian, blocks)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// func encodeBlockDbValue(w BlockType) []byte {
// 	value := make([]byte, 2)
// 	binary.LittleEndian.PutUint16(value, uint16(w))
// 	return value
// }

func decodeChunkDbValue(b []byte) []BlockType {
	if len(b) == 0 {
		log.Panic("len(b) == 0")
	}

	buff := bytes.NewBuffer(b)
	header, err := buff.ReadByte()
	if err != nil {
		log.Panicf("error. %s", err.Error())
	}

	if header == 0 {
		return nil
	}

	if len(b) != (ChunkWidth*ChunkWidth*ChunkWidth)*2+1 {
		log.Panicf("slice len[%d] is different from expected", len(b))
	}

	bts := make([]BlockType, ChunkWidth*ChunkWidth*ChunkWidth)
	binary.Read(buff, binary.LittleEndian, bts)

	return bts
}

func decodeBlockDbValue(b []byte) BlockType {
	if len(b) != 2 {
		log.Panicf("bad db value length:%d", len(b))
	}
	return BlockType(binary.LittleEndian.Uint16(b))
}
