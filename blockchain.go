package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Block represents an individual block in the blockchain.
type Block struct {
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	PoW          int
}

// Blockchain represents a chain of blocks.
type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

// calculateHash calculates and returns the hash of a block.
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.PoW)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

// mine mines a block to find a valid hash that satisfies the mining difficulty.
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.PoW++
		b.Hash = b.calculateHash()
	}
}

// CreateBlockchain creates a new blockchain with the specified mining difficulty.
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}
	return Blockchain{
		GenesisBlock: genesisBlock,
		Chain:        []Block{genesisBlock},
		Difficulty:   difficulty,
	}
}

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(from, to string, amount float64) error {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}

	previousBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Data:         blockData,
		PreviousHash: previousBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(bc.Difficulty)
	bc.Chain = append(bc.Chain, newBlock)

	// Validate the newly added block
	if !bc.IsValidBlock(len(bc.Chain) - 1) {
		return errors.New("invalid block")
	}

	return nil
}

// IsValidBlock checks the validity of a specific block in the blockchain.
func (bc *Blockchain) IsValidBlock(index int) bool {
	if index < 0 || index >= len(bc.Chain) {
		return false
	}

	block := bc.Chain[index]
	previousBlock := bc.Chain[index-1]

	// Check if block hash is valid
	if block.Hash != block.calculateHash() {
		return false
	}

	// Check if previous hash matches the hash of the previous block
	if block.PreviousHash != previousBlock.Hash {
		return false
	}

	return true
}

// IsValid checks the validity of the entire blockchain.
func (bc *Blockchain) IsValid() bool {
	// Skip the genesis block
	for i := 1; i < len(bc.Chain); i++ {
		if !bc.IsValidBlock(i) {
			return false
		}
	}
	return true
}

// DisplayBlockInfo formats the block's hash and associated information for printing.
func (b *Block) DisplayBlockInfo() string {
	info := fmt.Sprintf("Block Hash: %s\n", b.Hash)
	info += fmt.Sprintf("Previous Hash: %s\n", b.PreviousHash)
	info += fmt.Sprintf("Timestamp: %s\n", b.Timestamp.String())
	info += fmt.Sprintf("Proof-of-Work: %d\n", b.PoW)
	info += "Data:\n"
	for key, value := range b.Data {
		info += fmt.Sprintf("%s: %v\n", key, value)
	}
	return info
}

func main() {
	// Create a new blockchain instance with a mining difficulty of 5 (changeable)
	blockchain := CreateBlockchain(5)

	// Add blocks to the blockchain
	if err := blockchain.AddBlock("Cássio", "Luana", 4); err != nil {
		fmt.Println("Failed to add block:", err)
		return
	}

	if err := blockchain.AddBlock("Januário", "Luiz", 7); err != nil {
		fmt.Println("Failed to add block:", err)
		return
	}

	// Check the validity of the blockchain
	fmt.Println("Is blockchain valid?", blockchain.IsValid())

	// Display the hash and information of each block in the blockchain
	for i, block := range blockchain.Chain {
		fmt.Printf("Block %d:\n", i)
		fmt.Println(block.DisplayBlockInfo())
		fmt.Println()
	}
}
