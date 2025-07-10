package tgbotapp_test

import (
	"path/filepath"
	"reflect"
	"testing"

	tgbotapp "github.com/StridersTech2025/go-telegram-bot-app"
)

var (
	configTestFileYaml  = filepath.Join("testutil", "config_test.yaml")
	configTestFileJson  = filepath.Join("testutil", "config_test.json")
	configFileCorrupted = filepath.Join("testutil", "config_test_corrupt")
)

func TestConfigFromYamlShouldReadAndCreateAppConfig(t *testing.T) {
	// Arrange

	c := tgbotapp.NewAppConfig()

	// Act

	err := c.FromYaml(configTestFileYaml)

	// Assert
	if err != nil {
		t.Errorf("Should not return error. Found: %s", err.Error())
	}

	expected := tgbotapp.AppConfig{
		BotConfigs: []tgbotapp.BotConfig{
			tgbotapp.BotConfig{
				LanguageCode: "mm",
				Bot: tgbotapp.BotInfo{
					Name:             "အမေးအဖြေ",
					Description:      "အမေးအဖြေ bot",
					ShortDescription: "အတိုကောက်",
				},
			},
			tgbotapp.BotConfig{
				LanguageCode: "en",
				Bot: tgbotapp.BotInfo{
					Name:             "Quiz bot",
					Description:      "Quiz bot in yaml",
					ShortDescription: "Quiz bot short description yaml",
				},
			},
		},
	}

	if !reflect.DeepEqual(*c, expected) {
		t.Errorf("Expected AppConfig to be equal.\nExpected: %#v.\nFound: %#v", expected, *c)
	}

}

func TestConfigFromJsonShouldReadAndCreateAppConfig(t *testing.T) {
	// Arrange

	c := tgbotapp.NewAppConfig()

	// Act

	err := c.FromJson(configTestFileJson)

	// Assert

	if err != nil {
		t.Errorf("Should not return error.\nFound: %s", err.Error())
	}

	expected := tgbotapp.AppConfig{
		BotConfigs: []tgbotapp.BotConfig{
			tgbotapp.BotConfig{
				LanguageCode: "en",
				Bot: tgbotapp.BotInfo{
					Name:             "Quiz Bot 1",
					Description:      "Quiz bot long description",
					ShortDescription: "Quiz bot short description",
				},
			},
			tgbotapp.BotConfig{
				LanguageCode: "mm",
				Bot: tgbotapp.BotInfo{
					Name:             "အမေးအဖြေ",
					Description:      "အမေးအဖြေ bot",
					ShortDescription: "အမေးအဖြေ ဘော့ အတိုကောက်",
				},
			},
		},
	}

	if !reflect.DeepEqual(*c, expected) {
		t.Errorf("Expected AppConfig to be equal.\nExpected: %#v.\nFound: %#v", expected, *c)
	}

}

func TestConfigFromYamlShouldThrowErrorForReadingJson(t *testing.T) {
	// Arrange

	c := tgbotapp.NewAppConfig()

	// Act

	err := c.FromYaml(configTestFileJson)

	// Assert

	if err == nil {
		t.Errorf("Expects error to be returned.")
	} else {
		t.Logf("Got error: %s", err.Error())
	}

}

func TestConfigFromJsonShouldThrowErrorForReadingYaml(t *testing.T) {
	// Arrange

	c := tgbotapp.NewAppConfig()

	// Act

	err := c.FromJson(configTestFileYaml)

	// Assert

	if err == nil {
		t.Errorf("Expects error to be returned.")
	} else {
		t.Logf("Got error: %s", err.Error())
	}

}

func TestConfigFromYamlShouldThrowErrorForReadingInvalid(t *testing.T) {
	// Arrange

	c := tgbotapp.NewAppConfig()

	// Act

	err := c.FromYaml(configFileCorrupted)

	// Assert

	if err == nil {
		t.Errorf("Expects error to be returned.")
	} else {
		t.Logf("Got error: %s", err.Error())
	}

}

func TestConfigFromJsonShouldThrowErrorForReadingInvalid(t *testing.T) {
	// Arrange

	c := tgbotapp.NewAppConfig()

	// Act

	err := c.FromJson(configFileCorrupted)

	// Assert

	if err == nil {
		t.Errorf("Expects error to be returned.")
	} else {
		t.Logf("Got error: %s", err.Error())
	}

}
