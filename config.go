// Package envsubt provides YAML unmarshalling with environment variable substitution support.
package envsubt

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Unmarshal takes a YAML byte slice `in` and unmarshals it into the object `o`.
// It supports environment variable substitution for fields using the ${VAR_NAME} format.
func Unmarshal(in []byte, o any) error {
	ymlBytes, err := unmarshal(in)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(ymlBytes, o)
	if err != nil {
		return err
	}
	return nil
}

// validMapping checks if the input is a map[string]any and returns it.
// Returns an error if the input is not a valid mapping.
func validMapping(in any) (map[string]any, error) {
	comma, ok := in.(map[string]any)
	if !ok {
		err := fmt.Sprintf("Invalid config: got %T", in)
		return nil, NewConfigError(err)
	}
	return comma, nil
}

// unmarshal first converts the input YAML bytes into a map with resolved environment variables,
// then marshals it back into bytes for final unmarshalling.
func unmarshal(in []byte) ([]byte, error) {
	config, err := ymltoMap(in)
	if err != nil {
		return nil, err
	}
	newConfig, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	return newConfig, nil
}

// ymltoMap parses the YAML byte slice into a generic map and resolves any environment variables.
func ymltoMap(file []byte) (any, error) {
	var config any
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	err = resolveConfig(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// resolveConfig recursively walks through the config map and replaces environment variable placeholders.
func resolveConfig(config *any) error {
	MapConfig, err := validMapping(*config)
	if err != nil {
		return err
	}
	for k, v := range MapConfig {
		if MapConfig[k], err = resolveConfigVars(v); err != nil {
			return err
		}
	}
	return nil
}

// resolveConfigVars recursively resolves values in nested maps, replacing strings like ${VAR} with os.Getenv(VAR).
func resolveConfigVars(config any) (any, error) {
	MapConfig, err := validMapping(config)
	if err != nil {
		// Not a map, attempt to resolve as a single placeholder value.
		return resolvePlaceHolder(config), nil
	}
	for k, v := range MapConfig {
		if value, ok := v.(string); ok {
			MapConfig[k] = resolvePlaceHolder(value)
			continue
		}
		if MapConfig[k], err = resolveConfigVars(v); err != nil {
			return nil, err
		}
	}
	return config, nil // MapConfig is a reference to config
}

// resolvePlaceHolder checks if a string contains a placeholder of the form ${VAR}
// and replaces it with the value from the environment.
func resolvePlaceHolder(value any) any {
	strValue, ok := value.(string)
	if !ok {
		return value
	}
	if strings.Contains(strValue, "${") {
		last_index := len(strValue) - 1
		first_index := 2
		env_value := strValue[first_index:last_index]
		return os.Getenv(env_value)
	}
	return value
}

