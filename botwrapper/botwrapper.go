package botwrapper

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPIWrapper struct {
	*tgbotapi.BotAPI
}

func NewBotAPIWrapper(bot *tgbotapi.BotAPI) *BotAPIWrapper {
	return &BotAPIWrapper{bot}
}

func (b *BotAPIWrapper) SetConfig(c ConfigSetter) (bool, error) {

	params, err := c.params()

	if err != nil {
		return false, err
	}

	res, err := b.MakeRequest(c.method(), params)

	if err != nil {
		return false, err
	}

	return res.Ok, err

}

func (b *BotAPIWrapper) SendSetMyCommands(setMyCommands tgbotapi.SetMyCommandsConfig) (bool, error) {

	res, err := b.Request(setMyCommands)

	if err != nil {
		return false, err
	}

	return res.Ok, nil
}
