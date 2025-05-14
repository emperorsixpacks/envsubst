package envsubt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ClientSettings struct for testing
type ClientSettings struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Port    []int  `json:"port,omitempty" yaml:"port,omitempty"`
	Network string `json:"network,omitempty" yaml:"network,omitempty"`
}

func TestLoadConfig(t *testing.T) {
	// Test cases
	tests_c := struct {
		name        string
		yamlContent []byte
		setupEnv    map[string]string
		expected    ClientSettings
		expectError bool
		errorMsg    string
	}{
		name: "Valid YAML with placeholders",
		yamlContent: []byte(`
      name: ${CLIENT_NAME}
      version: v1.2.3
      port:
        - 30303
        - 8545
      network: testnet
      `),
		setupEnv: map[string]string{
			"CLIENT_NAME": "lighthouse",
		},
		expected: ClientSettings{
			Name:    "lighthouse",
			Version: "v1.2.3",
			Port:    []int{30303, 8545},
			Network: "testnet",
		},
		errorMsg: "",
	}

	t.Run(tests_c.name, func(t *testing.T) {
		// Setup file

		// Setup environment variables
		for k, v := range tests_c.setupEnv {
			os.Setenv(k, v)
			defer os.Unsetenv(k)
		}

		// Run LoadConfig
		var config ClientSettings
		err := Unmarshal(tests_c.yamlContent, &config)

		if err != nil {
			assert.Equal(t, tests_c.expected, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tests_c.expected, config)
		}
	})
}
