package tgbotapp

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type BotConfig struct {
	Name             string `yaml:"name" json:"name"`
	Description      string `yaml:"description" json:"description"`
	ShortDescription string `yaml:"shortDescription" json:"shortDescription"`
}

type AppConfig struct {
	LanguageCode string    `yaml:"languageCode" json:"languageCode"`
	Bot          BotConfig `yaml:"bot" json:"bot"`
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

// Read and construct configs from yaml file.
func (c *AppConfig) FromYaml(filename string) error {

	f, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading config file %q. %w", filename, err)
	}

	if isLikelyJson(f) {
		return fmt.Errorf("You might be parsing json file. use `FromJson` function.")
	}

	if c == nil {
		c = NewAppConfig()
	}

	if err = yaml.Unmarshal(f, c); err != nil {
		return fmt.Errorf("Error unmarshalling config file %q. %w", filename, err)
	}

	return nil
}

// Read and construct configs from json file.
func (c *AppConfig) FromJson(filename string) error {
	f, err := os.ReadFile(filename)

	if err != nil {
		return fmt.Errorf("Error reading config file %q. %w", filename, err)
	}

	if c == nil {
		c = NewAppConfig()
	}

	err = json.Unmarshal(f, c)

	if err != nil {
		return fmt.Errorf("Error unmarshalling config file %q. %w", filename, err)
	}

	return nil
}

func isLikelyJson(b []byte) bool {
	return b[0] == '{' || b[0] == '['
}
