package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leedinh/pluto/block"
	"github.com/leedinh/pluto/tracker"
)

func main() {
	client, err := ethclient.Dial("https://api-internal.roninchain.com/rpc")
	if err != nil {
		log.Fatal(err)
	}
	context := context.Background()

	blockSpan := block.NewBlockSpan(client, 7)
	blockTracker := tracker.NewBlockTracker(client, blockSpan)
	var wg sync.WaitGroup

	go blockSpan.BlockPoll(&context)
	fmt.Println("Start block tracker")
	go func() {
		for {
			select {
			case <-blockSpan.GetQueryChan():
				wg.Add(1)
				go blockTracker.Execute(&context, &wg)
			case <-context.Done():
				fmt.Println("Done")
				return
			}
			wg.Wait()
		}
	}()

	<-make(chan struct{})
}
