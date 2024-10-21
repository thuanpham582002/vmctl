package config

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".vmctl", "config.yaml")
}

func GetContextPath() (string, error) {
	file, err := os.Open(getConfigFilePath())
	if err != nil {
		return "", err
	}
	defer file.Close()

	config := &Config{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return "", err
	}

	return config.CurrentContext, nil
}

func GetContextDir() (string, error) {
	context, err := GetContextPath()
	if err != nil {
		return "", err
	}

	return filepath.Dir(context), nil

}

type Config struct {
	CurrentContext string `yaml:"current-context"`
}

func LoadConfigData() (map[string]interface{}, error) {
	context, err := GetContextPath()
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(context)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(file, &data); err != nil {
		// Fallback to JSON if YAML parsing fails
		if jsonErr := json.Unmarshal(file, &data); jsonErr != nil {
			return nil, jsonErr
		}
	}

	return data, nil
}
