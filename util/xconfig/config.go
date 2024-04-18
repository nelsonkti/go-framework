package xconfig

import (
	"encoding/json"
	"fmt"
)

const (
	Yaml = "yaml"
	Json = "json"
)

func New(c interface{}, reader ConfigReader) {
	rawConfig, err := reader.Load()
	if err != nil {
		panic(err)
	}

	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		panic(fmt.Errorf("failed to marshal config: %w", err))
	}

	if err := json.Unmarshal(configBytes, &c); err != nil {
		panic(fmt.Errorf("failed to unmarshal config into struct: %w", err))
	}
}
