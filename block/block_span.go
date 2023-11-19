package block

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type QueryEvent struct {
	From uint64
	To   uint64
}

type BlockSpan struct {
	client         *ethclient.Client
	confirm_blocks uint64
	QueryChan      chan uint64
}

func NewBlockSpan(client *ethclient.Client, confirm_blocks uint64) *BlockSpan {

	return &BlockSpan{
		client:         client,
		confirm_blocks: confirm_blocks,
		QueryChan:      make(chan uint64),
	}
}

func (bs *BlockSpan) BlockPoll(c *context.Context) {
	fmt.Println("Start block poll")

	for {
		lastestBlockNumber, err := bs.fetchLastestBlock()
		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Println("lastestBlockNumber: ", lastestBlockNumber)
		bs.QueryChan <- lastestBlockNumber - bs.confirm_blocks
		time.Sleep(3 * time.Second)
	}
}

func (bs *BlockSpan) fetchLastestBlock() (uint64, error) {
	header, err := bs.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return header.Number.Uint64(), nil
}
