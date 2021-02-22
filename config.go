package migrations

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config is used to map the migration config file.
type Config struct {
	Provider   string       `yaml:"provider"`
	Config     ConfigMap    `yaml:"config"`
	Migrations []*Migration `yaml:"migrations"`
}

// LoadConfigFromFile returns an instance of Config, populated
// with data from a YAML file.
func LoadConfigFromFile(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not read config file %s\n", filename)

		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("Config file does not contain valid YAML.\n")

		return nil, err
	}

	return &config, nil
}

// ConfigMap represents a map[string]interface{}, providing
// helper functions to access variables.
type ConfigMap map[string]interface{}

func (m ConfigMap) String(key string) (string, bool) {
	if m == nil {
		return "", false
	}

	v, ok := m[key]
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	if !ok {
		return "", false
	}

	return s, true
}
