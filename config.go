package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Resource string			`yaml:"resource"`
	RoutingKey string		`yaml:"routing_key"`
	CommandPre string		`yaml:"command_pre"`
	Command string			`yaml:"command"`
	CommandPost string		`yaml:"command_post"`
}

type Resources struct {
	Exchanges []Resource	`yaml:"exchanges"`
}

type Resource struct {
	Name string				`yaml:"name"`
	ResourceType string		`yaml:"type"`
	AutoDeleted bool		`yaml:"auto_deleted"`
	Internal bool			`yaml:"internal"`
	NoWait bool				`yaml:"no_wait"`
}

type Config struct {
	Resources Resources				`yaml:"resources"`
	Commands map[string]Command		`yaml:"commands"`
}

func parseConfigString(configContent *string) (Config, error) {
	config := Config{}

	err := yaml.Unmarshal([]byte(*configContent), &config)

	return config, err
}

func (container *Config) GetResource(name string) (Resource, error) {
	for _, item := range container.Resources.Exchanges {
		if item.Name == name {
			return item, nil
		}
	}

	return Resource{}, fmt.Errorf("Resource with name %s does not exist in configuration", name)
}

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