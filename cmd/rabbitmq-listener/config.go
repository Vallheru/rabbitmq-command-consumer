package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Command ...
type Command struct {
	Resource string			`yaml:"resource"`
	RoutingKey string		`yaml:"routing_key"`
	CommandPre string		`yaml:"command_pre"`
	Command string			`yaml:"command"`
	CommandPost string		`yaml:"command_post"`
}

// Resources ...
type Resources struct {
	Exchanges []Resource	`yaml:"exchanges"`
}

// Resource ...
type Resource struct {
	Category string
	Name string				`yaml:"name"`
	ResourceType string		`yaml:"type"`
	Durable bool			`yaml:"durable"`
	AutoDeleted bool		`yaml:"auto_deleted"`
	Internal bool			`yaml:"internal"`
	NoWait bool				`yaml:"no_wait"`
}

// ProgramConfig ...
type ProgramConfig struct {
	LogFilePath	string		`yaml:"log_file_path"`
}

// Config ...
type Config struct {
	Resources Resources				`yaml:"resources"`
	Commands map[string]Command		`yaml:"commands"`
	Program ProgramConfig			`yaml:"program"`
}

func parseConfigString(configContent *string) (Config, error) {
	config := Config{}

	err := yaml.Unmarshal([]byte(*configContent), &config)

	return config, err
}

// GetResource ...
func (container *Config) GetResource(name string) (Resource, error) {
	for _, item := range container.Resources.Exchanges {
		if item.Name == name {
			item.Category = ResourceExchange
			return item, nil
		}
	}

	return Resource{}, fmt.Errorf("Resource with name %s does not exist in configuration", name)
}

// GetCommand ...
func (container *Config) GetCommand(name string) (*Command, error) {
	item, ok := container.Commands[name];
	if !ok {
		return nil, fmt.Errorf("Command %s does not exist", name)
	}

	return &item, nil
}

// Validate ...
func (container *Config) Validate() (bool, error) {
	for key, item := range container.Commands {
		_, err := container.GetResource(item.Resource)

		if err != nil {
			return false, fmt.Errorf(
				"Resource %s is defined for %s command, but it is not declared",
				item.Resource,
				key,
			)
		}
	}

	return true, nil
}