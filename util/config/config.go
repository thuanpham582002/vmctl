package config

import (
	"encoding/json"
	"fmt"
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

func GetListContext() ([]string, error) {
	// Retrieve the data from Viper
	contexts := viper.Get("contexts")

	// Check if the data is of type []interface{}
	if contextList, ok := contexts.([]interface{}); ok {
		// Convert []interface{} to []string
		var result []string
		for _, item := range contextList {
			// Ensure each item is a string
			if str, ok := item.(string); ok {
				result = append(result, str)
			} else {
				// Return an error if any item is not a string
				return nil, fmt.Errorf("non-string item found in contexts")
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("contexts is not a list of strings")
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
