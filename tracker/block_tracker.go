package tracker

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leedinh/pluto/block"
	"github.com/leedinh/pluto/parser"
	"github.com/leedinh/pluto/utils"
)

const FIRST_BLOCK = 29430017 - 1000

type BlockTracker struct {
	id        string
	client    *ethclient.Client
	BlockSpan *block.BlockSpan
	cp        uint64
	db        *bolt.DB
}

func NewBlockTracker(client *ethclient.Client, db *bolt.DB, block_span *block.BlockSpan, tracker_id string) *BlockTracker {
	cp, err := LoadCheckPointForTracker(db, tracker_id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Start block tracker %s at block %d\n", tracker_id, cp)
	return &BlockTracker{
		id:        tracker_id,
		client:    client,
		db:        db,
		cp:        cp,
		BlockSpan: block_span,
	}
}

func LoadCheckPointForTracker(d *bolt.DB, id string) (uint64, error) {
	var cp uint64
	err := d.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tracker_check_points"))
		if b == nil {
			return fmt.Errorf("tracker bucket not found")
		}
		value := b.Get([]byte(id))
		if value == nil {
			return fmt.Errorf("checkpoint not found")
		}
		cp = binary.BigEndian.Uint64(value)

		return nil
	})
	if err != nil {
		return FIRST_BLOCK, err
	}
	return cp, nil
}

func SaveCheckPointForTracker(d *bolt.DB, id string, cp uint64) error {
	return d.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tracker_check_points"))
		if b == nil {
			return fmt.Errorf("tracker bucket not found")
		}
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, cp)
		return b.Put([]byte(id), buf)
	})
}

func (bt *BlockTracker) Execute(c *context.Context, wg *sync.WaitGroup, current_block uint64, rules []parser.Rule) {
	defer wg.Done()
	wg_block := sync.WaitGroup{}
	from := bt.cp + 1
	to := utils.Min(from+100, current_block)
	ch := make(chan struct{}, 10)
	if to-from < 1 {
		return
	}
	fmt.Printf("Execute block tracker %s from %d to %d\n", bt.id, from, to)
	for block_num := from; block_num <= to; block_num++ {
		ch <- struct{}{}
		wg_block.Add(1)
		go func(block_num uint64) {
			defer wg_block.Done()
			fmt.Println("Querying block ", block_num)
			block := parser.NewBlock(bt.client, c, block_num)
			parser.QueryTransactions(block.Transactions, rules)
			<-ch
		}(block_num)
	}
	wg_block.Wait()
	bt.cp = to
	err := SaveCheckPointForTracker(bt.db, bt.id, bt.cp)
	if err != nil {
		fmt.Println(err)
	}

}
