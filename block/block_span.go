package block

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/ethclient"
)

type QueryEvent struct {
	From uint64
	To   uint64
}

type BlockSpan struct {
	client             *ethclient.Client
	confirm_blocks     uint64
	LastestBlockNumber *uint64
}

func NewBlockSpan(client *ethclient.Client, confirm_blocks uint64) *BlockSpan {

	return &BlockSpan{
		client:             client,
		LastestBlockNumber: new(uint64),
	}
}

func (bs *BlockSpan) BlockPoll(c context.Context, wg *sync.WaitGroup) {
	fmt.Println("Start block poll")
	defer wg.Done()

	lastestBlockNumber, err := bs.fetchLastestBlock(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("lastestBlockNumber: ", lastestBlockNumber)
	atomic.StoreUint64(bs.LastestBlockNumber, lastestBlockNumber-bs.confirm_blocks)
}

func (bs *BlockSpan) fetchLastestBlock(c context.Context) (uint64, error) {
	header, err := bs.client.HeaderByNumber(c, nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return header.Number.Uint64(), nil
}
