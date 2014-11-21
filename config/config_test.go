package config_test

import (
	"github.com/kvisscher/cloudy-haystack/config"
	"reflect"
	"strings"
	"testing"
)

func TestConfigIsParsedCorrectly(t *testing.T) {
	jsonToParse := `
  {
    "targetBaseUrl": "http://wpos.platform.nl",
    "authToken": "some-very-random-token",
    "mappings": [
    {
      "from": "/update",
      "to": "/UpdateProductService"
    },
    {
      "from": "/recover",
      "to": "/RecoverService"
    }
    ]
  }
  `
	expectedMappingConfig := config.MappingConfig{
		TargetBaseUrl: "http://wpos.platform.nl",
		AuthToken:     "some-very-random-token",
		Mappings: []config.MappingPath{
			config.MappingPath{From: "/update", To: "/UpdateProductService"},
			config.MappingPath{From: "/recover", To: "/RecoverService"},
		},
	}

	mappings := config.Parse(strings.NewReader(jsonToParse))

	if !reflect.DeepEqual(mappings, expectedMappingConfig) {
		t.Errorf("expected %+v, got %+v", expectedMappingConfig, mappings)
	}
}
