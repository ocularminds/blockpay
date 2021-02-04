package main
import"ocularminds.com/blocks"

func main(){
	bc := blocks.NewBlockchain()
	/*
	bc.AddBlock("Send 1 BTC to Festus")
	bc.AddBlock("Send 2 BTC to Becki")

	for _, block := range bc.Blocks{
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
		pow := blocks.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}*/
	defer bc.Db.Close()
	cli := blocks.CLI{bc}
	cli.Run()
}

