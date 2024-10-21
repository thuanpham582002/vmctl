package model

import (
	"os"
	"vmctl/util/config"
	"vmctl/util/printcolor"
	"vmctl/util/yaml"
)

type Script struct {
	Command string        `yaml:"command"`
	Env     []Environment `yaml:"envs"`
}

func (s Script) GetCommand() (string, error) {
	dir, err := config.GetContextDir()
	if err != nil {
		return "", err
	}
	scriptText := ""
	// Add environment variables
	for _, env := range s.Env {
		environment, err := env.GetEnvironment()
		if err != nil {
			printcolor.Error(err.Error())
			continue
		}
		scriptText += environment
	}

	// Check if the command is a file
	if _, err := os.Stat(s.Command); err == nil {
		// Read the file
		content, err := os.ReadFile(s.Command)
		if err != nil {
			return "", err
		}
		scriptText += string(content)
	} else if _, err := os.Stat(dir + "/" + s.Command); err == nil {
		// Read the file
		content, err := os.ReadFile(dir + "/" + s.Command)
		if err != nil {
			return "", err
		}
		scriptText += string(content)
	} else {
		scriptText += s.Command
	}
	return scriptText, nil
}

type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	Kind  string `yaml:"type"`
}

func (e Environment) GetEnvironment() (string, error) {
	switch e.Kind {
	case "file":
		// Read file
		content, err := os.ReadFile(e.Value)
		if err != nil {
			return "", err
		}
		return "export " + e.Name + "=" + string(content) + ";", nil
	case "text":
		return "export " + e.Name + "=" + e.Value + ";", nil
	default:
		configData, err := config.LoadConfigData()
		if err != nil {
			return "", err
		}

		env, err := yaml.GetValueFromPath(configData, e.Value)
		if err != nil {
			return "", err
		}
		return "export " + e.Name + "=" + env.(string) + ";", nil
	}
}
