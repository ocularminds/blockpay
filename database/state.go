package database
import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct{
	Balances map[Account]uint
	txMempool []Tx
	dbFile *os.File
}

func NewStateFromDisk() (*State, error){
	//get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json");
	gen, err := loadGenesis(genFilePath)
	if err != nil {
      return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances{
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil{
		return nil, err
	}

	scanner := bufio.NewScanner(f);
	state := &State{balances, make([]Tx, 0), f}
	//Iterate over each the tx.db file's line
	for scanner.Scan(){
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		//Convert JSON encoded TX into object (struct)
		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)
		//Rebuild the state (user balances) as a series of events
		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}
	return state, nil
}

func (s *State) Add(tx Tx) error {
  if err := s.apply(tx); err != nil{
	  return err
  }
  s.txMempool = append(s.txMempool, tx)
  return nil
}

func (s *State) Persist() error{
  //make a copy of mempool because the s.txMempool will be modified
  //in the loop below
  mempool := make([]Tx, len(s.txMempool))
  copy(mempool, s.txMempool)
  for i := 0; i < len(mempool); i++ {
	  txJson, err := json.Marshal(s.txMempool[i])
	  if err != nil{
		  return err
	  }
	  if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
		  return err
	  }
	  //remove the Tx written to a file from the mempool
	  s.txMempool = append(s.txMempool[:i], s.txMempool[i+1:]...)
  }
  return nil;
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
	   s.Balances[tx.To] += tx.Value
	   return nil
	}
 
	if tx.Value > s.Balances[tx.From] {
	   return fmt.Errorf("insufficient balance")
	}
 
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
 
	return nil
 }