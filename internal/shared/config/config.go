package config

import (
	"bytes"
	"log"

	"github.com/spf13/viper"
)

type Variable struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type ProjectAction struct {
	Type     string `yaml:"type"`
	Template string `yaml:"template"`
	Output   string `yaml:"output"`
}
type ProjectSchematic struct {
	Description string          `yaml:"description"`
	Variables   []Variable      `yaml:"variables"`
	Actions     []ProjectAction `yaml:"actions"`
}

type PluginAction struct {
	Type       string `yaml:"type"`
	Template   string `yaml:"template"`
	Output     string `yaml:"output"`
	File       string `yaml:"file"`
	Import     string `yaml:"import"`
	Alias      string `yaml:"alias"`
	Dependency string `yaml:"dependency"`
	Route      string `yaml:"route"`
}
type PluginSchematic struct {
	Description string         `yaml:"description"`
	Variables   []Variable     `yaml:"variables"`
	Actions     []PluginAction `yaml:"actions"`
}

type RootSchematic struct {
	Module           string `yaml:"module"`
	GeneratorVersion string `yaml:"generatorVersion"`
}

func LoadSchematic[T any](configFileContent []byte) (T, error) {
	var config T

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(configFileContent))
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshalling config file: %v", err)
		return config, err
	}

	return config, nil
}
