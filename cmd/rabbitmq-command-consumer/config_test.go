package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var validData string = `
program:
  log_file_path: ./output.log
resources:
  exchanges:
    - name: sphinx_indexer_exchange
      type: direct
      durable: true
      auto_deleted: false
      internal: false
      no_wait: false
commands:
  product_variant_main:
    resource: sphinx_indexer_exchange
    routing_key: product_variant_main
    command_pre: echo "product_variant_main start"
    command: echo "product_variant_main command"
    command_post: echo "product_variant_main post"
  product_admin:
    resource: sphinx_indexer_exchange
    routing_key: product_admin
    command_pre: echo "product_admin start"
    command: echo "product_admin command"
    command_post: echo "product_admin post"
  product_variant_admin_inch:
    resource: sphinx_indexer_exchange
    routing_key: product_variant_admin_inch
    command_pre: echo "product_variant_admin_inch start"
    command: echo "product_variant_admin_inch command"
    command_post: echo "product_variant_admin_inch post"
`

var invalidData string = `
resources:
  exchanges:
    - name: new_sphinx_indexer_exchange
      type: direct
      auto_deleted: false
      internal: false
      no_wait: false
commands:
  product_variant_main:
    resource: sphinx_indexer_exchange
    routing_key: product_variant_main
    command_pre: echo "product_variant_main start"
    command: echo "product_variant_main command"
    command_post: echo "product_variant_main post"
`

func TestParseConfigStringValidData(t *testing.T) {

	result, err := parseConfigString(&validData)
	expected := Config{
		Program: ProgramConfig{
			LogFilePath: "./output.log",
		},
		Resources: Resources{
			Exchanges: []Resource{
				{
					Name: "sphinx_indexer_exchange",
					ResourceType: "direct",
					Durable: true,
					AutoDeleted: false,
					Internal: false,
					NoWait: false,
				},
			},
		},
		Commands: map[string]Command{
			"product_variant_main": Command{
				Resource: "sphinx_indexer_exchange",
				RoutingKey: "product_variant_main",
				CommandPre: "echo \"product_variant_main start\"",
				Command: "echo \"product_variant_main command\"",
				CommandPost: "echo \"product_variant_main post\"",
			},
			"product_admin": Command{
				Resource: "sphinx_indexer_exchange",
				RoutingKey: "product_admin",
				CommandPre: "echo \"product_admin start\"",
				Command: "echo \"product_admin command\"",
				CommandPost: "echo \"product_admin post\"",
			},
			"product_variant_admin_inch": Command{
				Resource: "sphinx_indexer_exchange",
				RoutingKey: "product_variant_admin_inch",
				CommandPre: "echo \"product_variant_admin_inch start\"",
				Command: "echo \"product_variant_admin_inch command\"",
				CommandPost: "echo \"product_variant_admin_inch post\"",
			},
		},
	}

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestInfrastructureValidation(t *testing.T) {
	testTable := []struct{
		data string
		isValid bool
		expectedError string
	}{
		{validData, 	true,  ""},
		{invalidData,	false, "sphinx_indexer_exchange is defined for product_variant_main"},
	}

	for _, item := range testTable {
		result, _ := parseConfigString(&item.data)

		validationResult, err := result.Validate()

		assert.Equal(t, item.isValid, validationResult)
		if item.expectedError != "" {
			assert.Contains(t, err.Error(), item.expectedError)
		}
	}
}

func TestGetResource(t *testing.T) {
	testTable := []struct {
		data string
		wantedResource string
		expectedRes Resource
		expectedError *string
	}{
		{
			validData,
			"sphinx_indexer_exchange",
			Resource{
				Category: ResourceExchange,
				Name: "sphinx_indexer_exchange",
				ResourceType: "direct",
				AutoDeleted: false,
				Durable: true,
				Internal: false,
				NoWait: false,
			},
			nil,
		},
		{
			validData,
			"unknown_resource",
			Resource{},
			&(&struct{ x string }{"unknown_resource does not exist",}).x,
		},
	}

	for _, item := range testTable {
		result, _ := parseConfigString(&item.data)
		resource, err := result.GetResource(item.wantedResource)
		assert.Equal(t, item.expectedRes, resource)
		if item.expectedError == nil {
			assert.Nil(t, err)
		} else {
			assert.Contains(t, err.Error(), *item.expectedError)
		}
		
	}
}

func TestParseConfigStringInvalidData(t *testing.T) {
	data := `
	test:
	test1:
    `

	result, err := parseConfigString(&data)

	expected := Config{}
	assert.Equal(t, expected, result)
	assert.Contains(t, err.Error(), "yaml: line 2: found character that cannot start any token")
}
