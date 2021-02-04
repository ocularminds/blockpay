package blocks
import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
	"fmt"
	"encoding/gob"
)

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

func NewGenesisBlock() *Block{
	return NewBlock("Genesis Block", []byte{})
}

func (b *Block) Serialize() []byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)  
	err := encoder.Encode(b)
	fmt.Println(err)
	return result.Bytes()
}