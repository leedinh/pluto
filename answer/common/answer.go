package common

import "github.com/leedinh/pluto/model"

func TrackerMessage() model.Message {
	return model.Message{
		Text: "New event alert",
	}
}

func GreetMessage() model.Message {
	return model.Message{
		Text: "Welcome to the Tracker Bot",
	}
}

func NotPermission() model.Message {
	return model.Message{
		Text: "You do not have permission to use this command",
	}
}

func Success() model.Message {
	return model.Message{
		Text: "Success",
	}
}

func Fail() model.Message {
	return model.Message{
		Text: "Fail",
	}
}

func CustomMessage(text string) model.Message {
	return model.Message{
		Text: text,
	}
}
