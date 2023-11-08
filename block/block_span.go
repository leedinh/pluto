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
	From           uint64
	To             uint64
	confirm_blocks uint64
	query_chan     chan QueryEvent
}

func NewBlockSpan(client *ethclient.Client, confirm_blocks uint64) *BlockSpan {
	return &BlockSpan{client: client, confirm_blocks: confirm_blocks, query_chan: make(chan QueryEvent)}
}

func (bs *BlockSpan) GetQueryChan() chan QueryEvent {
	return bs.query_chan
}

func (bs *BlockSpan) BlockPoll(c *context.Context) {
	fmt.Println("Start block poll")

	for {
		header, err := bs.client.HeaderByNumber(*c, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		lastestBlockNumber := header.Number.Uint64()
		fmt.Println("lastestBlockNumber: ", lastestBlockNumber)
		bs.refresh(lastestBlockNumber)
		time.Sleep(5 * time.Second)
	}
}

func (bs *BlockSpan) GetLog() string {
	return fmt.Sprintf("from: %d to: %d", bs.From, bs.To)
}

func (bs *BlockSpan) refresh(lastestBlockNumber uint64) {
	if bs.From == 0 {
		bs.To = lastestBlockNumber - bs.confirm_blocks
		bs.From = bs.To - 200
		bs.query_chan <- QueryEvent{From: bs.From, To: bs.To}
	} else {
		if bs.To < lastestBlockNumber-bs.confirm_blocks {
			if (bs.To+1)-(lastestBlockNumber-bs.confirm_blocks) < 1 {
				return
			}
			bs.From = bs.To + 1
			bs.To = lastestBlockNumber - bs.confirm_blocks
			bs.query_chan <- QueryEvent{From: bs.From, To: bs.To}
			return
		}
	}
}
