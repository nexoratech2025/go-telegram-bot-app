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
		LanguageCode: "mm",
		Bot: tgbotapp.BotConfig{
			Name:             "အမေးအဖြေ",
			ShortDescription: "အတိုကောက်",
			Description:      "အမေးအဖြေ bot",
		},
	}

	if !reflect.DeepEqual(*c, expected) {
		t.Errorf("Expected AppConfig to be equal. Expected: %#v, Found: %#v", expected, *c)
	}

}

func TestConfigFromJsonShouldReadAndCreateAppConfig(t *testing.T) {
	// Arrange

	c := tgbotapp.NewAppConfig()

	// Act

	err := c.FromJson(configTestFileJson)

	// Assert

	if err != nil {
		t.Errorf("Should not return error. Found: %s", err.Error())
	}

	expected := tgbotapp.AppConfig{
		LanguageCode: "en",
		Bot: tgbotapp.BotConfig{
			Name:             "Quiz Bot 1",
			ShortDescription: "Quiz bot short description",
			Description:      "Quiz bot long description",
		},
	}

	if !reflect.DeepEqual(*c, expected) {
		t.Errorf("Expected AppConfig to be equal. Expected: %#v, Found: %#v", expected, *c)
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
