package rpc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)


//Block 定义数据类型
type Block struct {
	Index         int64  //区块编号
	Timestamp     int64  //区块时间戳
	PrevBlockHash string //上一个区块的哈希值
	Hash          string //当前区块的哈希值

	Data string //区块数据
}
//calculateHash 计算哈希值
func CalculateHash(b Block) string {
	blockData := string(b.Index) + string(b.Timestamp) + b.PrevBlockHash + b.Data
	hashInBytes := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hashInBytes[:])
}

//GenerateNewBlock 生成区块
func GenerateNewBlock(preBlock Block, data string) Block {
	newBlock := Block{}
	newBlock.Index = preBlock.Index + 1
	newBlock.Timestamp = time.Now().Unix()
	newBlock.PrevBlockHash = preBlock.Hash
	newBlock.Data = data
	newBlock.Hash = CalculateHash(newBlock)
	return newBlock
}

//GenerateGenesisBlock 第一个区块
func GenerateGenesisBlock() Block {
	preBlock := Block{}
	preBlock.Index = -1
	preBlock.Hash = ""
	return GenerateNewBlock(preBlock, "Genesis Block")
}




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
		fmt.Println()
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

var blockchain *Blockchain

func run(){
	http.HandleFunc("/blockchain/get",blockchainGetHandler)
	http.HandleFunc("/blockchain/write",blockchainWriteHandler)
	http.ListenAndServe("localhost:4000",nil)
}

func blockchainGetHandler(w http.ResponseWriter,r *http.Request){
	bytes,error := json.Marshal(blockchain)
	if error != nil {
		http.Error(w,error.Error(),http.StatusInternalServerError)
		return
	}
	io.WriteString(w,string(bytes))
}

func blockchainWriteHandler(w http.ResponseWriter,r *http.Request)  {
	blockData := r.URL.Query().Get("data")
	blockchain.SendData(blockData)
	blockchainGetHandler(w,r)
}

func main(){
	blockchain = NewBlockchain()
	run()
}