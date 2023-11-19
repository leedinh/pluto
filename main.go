package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/leedinh/pluto/block"
	"github.com/leedinh/pluto/bot"
	"github.com/leedinh/pluto/db"
	"github.com/leedinh/pluto/logger"
	"github.com/leedinh/pluto/model"
	"github.com/leedinh/pluto/parser"
	"github.com/leedinh/pluto/sc"
	"github.com/leedinh/pluto/tracker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	logger := logger.GetLogger()

	fmt.Println("Start Pluto")
	db, close, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	wl := sc.NewSMWhitelist(db)
	log.Println(wl.Whitelist)

	rpc_client, err := ethclient.Dial("https://api-internal.roninchain.com/rpc")
	if err != nil {
		log.Fatal(err)
	}
	defer rpc_client.Close()
	context := context.Background()

	blockSpan := block.NewBlockSpan(rpc_client, 7)
	blockTracker := tracker.NewBlockTracker(rpc_client, db.Db, blockSpan, "tracker_1")

	// tx := parser.QueryTransactionByHash(rpc_client, &context, "0xf5669b3d863db1a70bfde8f7d0fc148152a28aaf32e2bce55b2683bf9a020643")
	// log.Println(tx.To)
	// log.Println(sc.GetContractABI(strings.ToLower(tx.To)))
	var wg sync.WaitGroup

	go blockSpan.BlockPoll(&context)
	fmt.Println("Start block tracker")
	rules := parser.InitRules()
	go func() {
		for {
			select {
			case comming_block := <-blockSpan.QueryChan:
				wg.Add(1)
				go blockTracker.Execute(&context, &wg, comming_block, rules)
			case <-context.Done():
				fmt.Println("Done")
				return
			}
			wg.Wait()
		}
	}()

	bot := bot.InitBot(telegramBotToken, &logger, model.Flow{}, db)
	bot.Start()
	// <-make(chan struct{})
}
