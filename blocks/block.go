package blocks
import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
)

type Block struct{
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce int
}

type Blockchain struct {
	Blocks []*Block
}

func (b *Block) SetHash(){
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block{
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	//block.SetHash()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (bc *Blockchain) AddBlock(data string){
	prevBlock := bc.Blocks[len(bc.Blocks) - 1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock() *Block{
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain{
    return &Blockchain{[]*Block{NewGenesisBlock()}}
}