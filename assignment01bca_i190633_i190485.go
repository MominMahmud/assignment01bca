package main

import (
	"fmt"
	"log"
	"encoding/binary"
	"math"
	"math/big"
	"strconv"
	"bytes"
	"crypto/sha256"
	
	
)

const diff = 10 

type AsingleBlock struct {
	Blocks []*block 
}

type block struct {
	Hash     []byte 
	Data     []byte
	PrevHash []byte
	Nonce    int
}

type proofWork struct {
	block  *block
	Target *big.Int
}
type MtreeMaker struct {
	RootNode *Mnode
}

type Mnode struct {
	Left    *Mnode
	Right   *Mnode
	Data    []byte
	endNode    bool
	Content []byte
}

func CreateBlockChain(array [][]byte, prevHash []byte) *block {
	tree := NewMerkleTree(array)

	block := &block{[]byte{}, []byte(tree.RootNode.Data), prevHash, 0}
	pow := setProofofWork(block)
	nonce, hash := pow.MineBlock()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (chain *AsingleBlock) addBlock(array [][]byte) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlockChain(array, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

func returnGenesisBlock(array [][]byte) *block {
	return CreateBlockChain(array, []byte{})
}

func initiateChain(array [][]byte) *AsingleBlock {
	return &AsingleBlock{[]*block{returnGenesisBlock(array)}}
}

func setProofofWork(b *block) *proofWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-diff))
	pow := &proofWork{b, target}
	return pow
}

func (pow *proofWork) initDataBlock(nonce int) []byte {
	data := bytes.Join([][]byte{pow.block.PrevHash, pow.block.Data, hexConvert(int64(nonce)), hexConvert(int64(diff))}, []byte{})
	return data
}

func (pow *proofWork) MineBlock() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.initDataBlock(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *proofWork) VerifyChange() bool {
	var intHash big.Int
	data := pow.initDataBlock(pow.block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

func hexConvert(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func DisplayBlocks(chain *AsingleBlock) {
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := setProofofWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.VerifyChange()))
		fmt.Println()
	}
}

func NewMerkleNode(left, right *Mnode, data []byte) *Mnode {
	node := Mnode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
		node.endNode = true
		node.Content = data
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
		node.endNode = false
	}

	node.Left = left
	node.Right = right

	return &node
}

func NewMerkleTree(data [][]byte) *MtreeMaker {
	var nodes []Mnode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, *node)
	}

	for i := 0; i < len(data)/2; i++ {
		var level []Mnode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	tree := MtreeMaker{&nodes[0]}

	return &tree
}

func main() {

	var array [][]byte
	array = append(array, []byte("Genesis data 1"))
	array = append(array, []byte("Genesis data 2"))
	array = append(array, []byte("Genesis data 3"))
	//fmt.Printf("%s\n", array[1])
	chain := initiateChain(array)
	var array1 [][]byte
	array1 = append(array1, []byte("First block after Genesis data 1"))
	array1 = append(array1, []byte("First block after Genesis data 2"))
	array1 = append(array1, []byte("First block after Genesis data 3"))
	array1 = append(array1, []byte("First block after Genesis data 4"))
	chain.addBlock(array1)
	var array2 [][]byte
	array2 = append(array2, []byte("First block after Genesis data 1"))
	array2 = append(array2, []byte("First block after Genesis data 2"))
	array2 = append(array2, []byte("First block after Genesis data 3"))
	chain.addBlock(array2)
	//DisplayBlocks(chain)

}