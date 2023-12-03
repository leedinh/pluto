package parser

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leedinh/pluto/model"
)

type Block struct {
	block        *types.Block
	Number       int64
	Hash         string
	Time         uint64
	Transactions *[]Transaction
}

type Transaction struct {
	Hash  string
	To    string
	Value string
	Nonce uint64
	Data  []byte
}

func GetTransactions(block *types.Block) *[]Transaction {
	var transactions []Transaction
	for _, tx := range block.Transactions() {
		transactions = append(transactions, Transaction{
			Hash:  tx.Hash().Hex(),
			To:    tx.To().Hex(),
			Value: tx.Value().String(),
			Nonce: tx.Nonce(),
			Data:  tx.Data(),
		})
	}
	return &transactions
}

func NewBlock(client *ethclient.Client, c *context.Context, number uint64) *Block {
	blockNumber := big.NewInt(int64(number))
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	return &Block{
		block:        block,
		Number:       block.Number().Int64(),
		Hash:         block.Hash().Hex(),
		Time:         block.Time(),
		Transactions: GetTransactions(block),
	}
}

func QueryTransactions(txs *[]Transaction, rule []Rule, tupdate *model.TrackerUpdate) {
	for _, tx := range *txs {
		for _, r := range rule {
			if r.Check(&tx) {
				log.Println("Found transaction ", tx.Hash)
				tupdate.Ch <- model.TrackerEvent{
					EventType: r.Name,
				}
			}
		}
	}
}

func QueryTransactionByHash(client *ethclient.Client, c *context.Context, hash string) Transaction {
	tx, _, err := client.TransactionByHash(*c, common.HexToHash(hash))
	if err != nil {
		log.Fatal(err)
	}
	return Transaction{
		Hash:  tx.Hash().Hex(),
		To:    tx.To().Hex(),
		Value: tx.Value().String(),
		Nonce: tx.Nonce(),
		Data:  tx.Data(),
	}
}
