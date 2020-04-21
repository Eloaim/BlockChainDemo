package core

import (
"crypto/sha256"
"encoding/hex"
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


