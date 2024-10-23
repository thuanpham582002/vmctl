package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func GetContextPath() (string, error) {
	if path, ok := viper.Get("current-context").(string); ok {
		return path, nil
	}
	return "", nil
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
