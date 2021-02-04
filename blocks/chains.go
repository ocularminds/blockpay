package blocks

import (
	"fmt"
	"github.com/boltdb/bolt"
	"bytes"
	"encoding/gob"
	"log"
) 

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

/*Open a DB file.
Check if there’s a blockchain stored in it.
If there’s a blockchain:
Create a new Blockchain instance.
Set the tip of the Blockchain instance to the last block hash stored in the DB.
If there’s no existing blockchain:
Create the genesis block.
Store in the DB.
Save the genesis block’s hash as the last block hash.
Create a new Blockchain instance with its tip pointing at the genesis block.
*/
func NewBlockchain() *Blockchain{
	//return &Blockchain{[]*Block{NewGenesisBlock()}}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil{
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("1"), genesis.Hash)
			tip = genesis.Hash
		}else{
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}

func (bc *Blockchain) AddBlock(data string){
	var lastHash []byte
	//prevBlock := bc.Blocks[len(bc.Blocks) - 1]
	//newBlock := NewBlock(data, prevBlock.Hash)
	//bc.Blocks = append(bc.Blocks, newBlock)
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
        return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)
	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.Db}

	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = Deserialize(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}

func Deserialize(data []byte) *Block{
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	fmt.Println(err)
	return &block
  }
