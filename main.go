package main

import (
	"context"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leedinh/pluto/parser"
	"github.com/leedinh/pluto/sc"
)

func main() {
	rpc_client, err := ethclient.Dial("https://api-internal.roninchain.com/rpc")
	if err != nil {
		log.Fatal(err)
	}
	defer rpc_client.Close()

	db, close, err := sc.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	wl := sc.NewSMWhitelist(db)
	log.Println(wl.Whitelist)

	context := context.Background()

	// blockSpan := block.NewBlockSpan(client, 7)
	// blockTracker := tracker.NewBlockTracker(client, blockSpan)

	tx := parser.QueryTransactionByHash(rpc_client, &context, "0xf5669b3d863db1a70bfde8f7d0fc148152a28aaf32e2bce55b2683bf9a020643")
	log.Println(tx.To)
	log.Println(sc.GetContractABI(strings.ToLower(tx.To)))
	// var wg sync.WaitGroup

	// go blockSpan.BlockPoll(&context)
	// fmt.Println("Start block tracker")
	// go func() {
	// 	for {
	// 		select {
	// 		case <-blockSpan.GetQueryChan():
	// 			wg.Add(1)
	// 			go blockTracker.Execute(&context, &wg)
	// 		case <-context.Done():
	// 			fmt.Println("Done")
	// 			return
	// 		}
	// 		wg.Wait()
	// 	}
	// }()

	// <-make(chan struct{})
}
