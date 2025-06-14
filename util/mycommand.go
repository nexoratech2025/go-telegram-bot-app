package util

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendSetMyCommands(botAPI tgbotapi.BotAPI, setMyCommands tgbotapi.SetMyCommandsConfig) (bool, error) {

	res, err := botAPI.Request(setMyCommands)

	if err != nil {
		return false, err
	}

	return res.Ok, nil
}
