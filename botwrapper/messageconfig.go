package botwrapper

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ConfigSetter interface {
	method() string
	params() (tgbotapi.Params, error)
}

type SetMyNameConfig struct {
	Name         string
	LanguageCode string
}

func NewSetMyName(name string, languageCode string) SetMyNameConfig {
	return SetMyNameConfig{
		Name:         name,
		LanguageCode: languageCode,
	}
}

func (c SetMyNameConfig) method() string {
	return "setMyName"
}

func (c SetMyNameConfig) params() (tgbotapi.Params, error) {
	params := make(tgbotapi.Params)

	params.AddNonEmpty("name", c.Name)
	params.AddNonEmpty("language_code", c.LanguageCode)

	return params, nil
}

type SetMyDescription struct {
	Description  string
	LanguageCode string
}

func NewSetMyDescription(name string, languageCode string) SetMyDescription {
	return SetMyDescription{
		Description:  name,
		LanguageCode: languageCode,
	}
}

func (c SetMyDescription) method() string {
	return "setMyDescription"
}

func (c SetMyDescription) params() (tgbotapi.Params, error) {
	params := make(tgbotapi.Params)

	params.AddNonEmpty("description", c.Description)
	params.AddNonEmpty("language_code", c.LanguageCode)

	return params, nil
}

type SetMyShortDescription struct {
	ShortDescription string
	LanguageCode     string
}

func NewSetMyShortDescription(name string, languageCode string) SetMyShortDescription {
	return SetMyShortDescription{
		ShortDescription: name,
		LanguageCode:     languageCode,
	}
}

func (c SetMyShortDescription) method() string {
	return "setMyShortDescription"
}

func (c SetMyShortDescription) params() (tgbotapi.Params, error) {
	params := make(tgbotapi.Params)

	params.AddNonEmpty("short_description", c.ShortDescription)
	params.AddNonEmpty("language_code", c.LanguageCode)

	return params, nil
}
