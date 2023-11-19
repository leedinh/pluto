package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/leedinh/pluto/answer/common"
	"github.com/leedinh/pluto/answer/profile"
	"github.com/leedinh/pluto/model"
	"go.uber.org/zap"
)

func (b *Bot) UpdateRouter(upd tgbotapi.Update) {
	updLocal := model.DecodeToLocal(upd)
	if msg := upd.Message; msg != nil {
		if msg.IsCommand() {
			b.SendMessage(b.CommandsHandler(upd.Message.Command(), updLocal))
		} else {
			b.SendMessage(b.MessageHandler(upd))
		}
	}
	if cq := upd.CallbackQuery; cq != nil {
		b.SendMessage(b.CallbacksHandler(updLocal))
	}
}

func (b *Bot) CommandsHandler(command string, updLocal *model.UpdateLocal) tgbotapi.Chattable {
	/*
		your commands processing logic should be here
		return <message>
	*/
	switch command {
	case "start":
		return common.GreetMessage().BuildBotMessage(int64(updLocal.TelegramChatID))
	case "profile":
		return profile.Create().BuildBotMessage(int64(updLocal.TelegramChatID))
	case "add_sc":
		if int64(updLocal.TelegramUserID) != 2100956682 {
			common.NotPermission().BuildBotMessage(int64(updLocal.TelegramChatID))
		}
		if err := model.AddSCToWhitelist(b.Db, updLocal.Message); err != nil {
			b.Logger.Error("error", zap.String("reason", err.Error()))
			return common.Fail().BuildBotMessage(int64(updLocal.TelegramChatID))
		}
		return common.Success().BuildBotMessage(int64(updLocal.TelegramChatID))

	case "list_sc":
		if int64(updLocal.TelegramUserID) != 2100956682 {
			common.NotPermission().BuildBotMessage(int64(updLocal.TelegramChatID))
		}
		sc_map, err := model.GetListSC(b.Db)
		if err != nil {
			b.Logger.Error("error", zap.String("reason", err.Error()))
			return common.Fail().BuildBotMessage(int64(updLocal.TelegramChatID))
		}
		return common.CustomMessage(ParseListSC(sc_map)).BuildBotMessage(int64(updLocal.TelegramChatID))
	}

	return common.TrackerMessage().BuildBotMessage(int64(updLocal.TelegramChatID))
}

func (b *Bot) MessageHandler(upd tgbotapi.Update) tgbotapi.Chattable {
	updLocal := model.DecodeToLocal(upd)
	/*
		your message processing logic should be here
		return <message>
	*/
	return common.TrackerMessage().BuildBotMessage(int64(updLocal.TelegramChatID))
}

func (b *Bot) CallbacksHandler(updLocal *model.UpdateLocal) tgbotapi.Chattable {
	cData := updLocal.CallbackData
	replyMessage, err := b.Flow.Handle(&cData, updLocal)
	if err != nil {
		b.Logger.Error("error", zap.String("reason", err.Error()))
		return common.TrackerMessage().BuildBotMessage(int64(updLocal.TelegramChatID))
	}
	return replyMessage
}

func ParseListSC(sc_map map[string]model.ContractInfo) string {
	var list string
	for k, v := range sc_map {
		list += k + " - " + v.Name + "\n"
	}
	return list
}
