package envsubt

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

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

func validMapping(in any) (map[string]any, error) {
	comma, ok := in.(map[string]any)
	if !ok {
		err := fmt.Sprintf("Invalid condfig: got %s", comma)
		return nil, NewConfigError(err)
	}
	return comma, nil
}

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

func resolveConfigVars(config any) (any, error) {
	MapConfig, err := validMapping(config)
	if err != nil {
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
