package core
import "log"
import "fmt"

type Blockchain struct {
	Blocks []*Block
}
func (bc *Blockchain) Print(){
	for _,block := range bc.Blocks{
		fmt.Printf("Index: %d\n",block.Index)
		fmt.Printf("Prev: %s\n",block.PrevBlockHash)
		fmt.Printf("Curr: %s\n",block.Hash)
		fmt.Printf("Data: %s\n",block.Data)
		fmt.Printf("Timestamp: %d\n",block.Timestamp)
	}
}
func NewBlockchain() *Blockchain{
	genesisBlock := GenerateGenesisBlock()
	blockchain := Blockchain{}
	blockchain.ApendBlock(&genesisBlock)
	return &blockchain
}

func (bc *Blockchain)SendData(data string){
	preBlock := *bc.Blocks[len(bc.Blocks)-1]
	newBlock := GenerateNewBlock(preBlock,data)
	bc.ApendBlock(&newBlock)
}

func (bc *Blockchain) ApendBlock(newBlock *Block) {
	if len(bc.Blocks) == 0{
		bc.Blocks = append(bc.Blocks,newBlock)
		return
	}
	if IsVald(*newBlock, *bc.Blocks[len(bc.Blocks)-1]) {
		bc.Blocks = append(bc.Blocks, newBlock)
	} else {
		log.Fatal("invail block")
	}

}

func IsVald(newBlock Block, oldBlock Block) bool {
	if newBlock.Index != oldBlock.Index+1 {
		return false
	}
	if newBlock.PrevBlockHash != oldBlock.Hash {
		return false
	}
	if newBlock.Hash != CalculateHash(newBlock) {
		return false
	}
	return true
}