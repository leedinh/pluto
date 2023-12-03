package main

import (
	"context"
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
	"github.com/leedinh/pluto/tracker"
	"go.uber.org/zap"
)

var (
	telegramBotToken string
	logg             *zap.Logger
	d                *db.Database
	rpcClient        *ethclient.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	logg = logger.GetLogger()
	logg.Info("Start Pluto")

	d, _, err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("Connect to DB successfully", zap.Any("database", d))

	rpcClient, err = ethclient.Dial("https://api-internal.roninchain.com/rpc")
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("Connect to RPC successfully", zap.Any("rpcClient", rpcClient))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var trackerQueue []interface{}

	blockSpan := block.NewBlockSpan(rpcClient, 7)
	eventTracker := tracker.NewBlockTracker(rpcClient, d.Db, blockSpan, "tracker_1")
	trackerQueue = append(trackerQueue, blockSpan, eventTracker)
	trackerUpdate := model.NewTrackerUpdate()
	var trackerWG sync.WaitGroup
	rules := parser.InitRules()

	go func() {
		for {
			for _, item := range trackerQueue {
				trackerWG.Add(1)
				switch it := item.(type) {
				case *block.BlockSpan:
					go it.BlockPoll(ctx, &trackerWG)
				case *tracker.BlockTracker:
					go it.Execute(ctx, &trackerWG, it.BlockSpan.LastestBlockNumber, rules, trackerUpdate)
				}
			}
			trackerWG.Wait()
		}
	}()

	bot := bot.InitBot(telegramBotToken, logg, model.Flow{}, d, trackerUpdate)
	bot.Start()
}
