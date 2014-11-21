package config

import (
	"encoding/json"
	"io"
)

type MappingPath struct {
	From string
	To   string
}

type MappingConfig struct {
	TargetBaseUrl string        `json:"targetBaseUrl"`
	AuthToken     string        `json:"authToken"`
	Mappings      []MappingPath `json:"mappings"`
}

func Parse(reader io.Reader) (config MappingConfig) {
	decoder := json.NewDecoder(reader)

	decoder.Decode(&config)

	return config
}
