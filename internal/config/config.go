package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := home + "/" + configFileName

	return path, nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	gatorConfig := Config{}
	if err := json.Unmarshal(fileContent, &gatorConfig); err != nil {
		return Config{}, err
	}

	return gatorConfig, nil
}

func (c *Config) SetUser(user_name string) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	c.CurrentUserName = user_name

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
