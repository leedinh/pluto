package tracker

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leedinh/pluto/block"
	"github.com/leedinh/pluto/parser"
)

type BlockTracker struct {
	client     *ethclient.Client
	block_span *block.BlockSpan
}

func NewBlockTracker(client *ethclient.Client, block_span *block.BlockSpan) *BlockTracker {
	return &BlockTracker{client, block_span}
}

func (bt *BlockTracker) Start() {
	fmt.Println("Start block tracker")
	for {
		fmt.Println(bt.block_span.GetLog())
	}
}

func (bt *BlockTracker) Execute(c *context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	bs := bt.GetBlockSpan()
	fmt.Println(bs.GetLog())
	wg_block := sync.WaitGroup{}
	for block_num := bs.From; block_num <= bs.To; block_num++ {
		wg_block.Add(1)
		go func(block_num uint64) {
			defer wg_block.Done()
			block := parser.NewBlock(bt.client, c, block_num)
			fmt.Println("Querying block ", block.Number)
			parser.QueryTransactions(block.Transactions)
		}(block_num)
	}
	wg_block.Wait()

}

func (bt *BlockTracker) GetBlockSpan() *block.BlockSpan {
	return bt.block_span
}
