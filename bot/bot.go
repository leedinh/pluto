package bot

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/leedinh/pluto/db"
	"github.com/leedinh/pluto/model"
	"go.uber.org/zap"
)

type Bot struct {
	API      *tgbotapi.BotAPI
	commands map[string]func()
	Logger   *zap.Logger
	Flow     model.Flow
	Db       *bolt.DB
}

func InitBot(token string, logger *zap.Logger, flow model.Flow, db *db.Database) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	return &Bot{
		API:      bot,
		commands: make(map[string]func()),
		Logger:   logger,
		Flow:     flow,
		Db:       db.Db,
	}
}

func (b *Bot) Start() {
	fmt.Println("Bot started")
	defer func() {
		fmt.Println("Bot stopped")
	}()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		go b.UpdateRouter(update)
	}
}

func (b *Bot) SendMessage(msg tgbotapi.Chattable) {
	_, err := b.API.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
