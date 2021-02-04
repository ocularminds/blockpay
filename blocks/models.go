package blocks

import "github.com/boltdb/bolt"
type Block struct{
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce int
}

type Blockchain struct {
	//Blocks []*Block
	tip []byte
	Db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

type CLI struct {
	Chains *Blockchain
}