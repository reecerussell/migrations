package migrations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile_GivenValidFilename_ReturnsConfig(t *testing.T) {
	const testYAML = `provider: test
config:
  value: Hello World
migrations:
- name: Test
  up: test.up.sql
  down: test.down.sql`

	file, err := os.Create("TestLoadConfigFromFile_GivenValidFilename_ReturnsConfig")
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
		return
	}

	file.Write([]byte(testYAML))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestLoadConfigFromFile_GivenValidFilename_ReturnsConfig")
	})

	conf, err := LoadConfigFromFile("TestLoadConfigFromFile_GivenValidFilename_ReturnsConfig")
	assert.NoError(t, err)
	assert.Equal(t, "test", conf.Provider)
	assert.Equal(t, "Hello World", conf.Config["value"])
	assert.Equal(t, 1, len(conf.Migrations))
	assert.Equal(t, "Test", conf.Migrations[0].Name)
	assert.Equal(t, "test.up.sql", conf.Migrations[0].UpFile)
	assert.Equal(t, "test.down.sql", conf.Migrations[0].DownFile)
}

func TestLoadConfigFromFile_GivenNonExistantFilename_ReturnsError(t *testing.T) {
	conf, err := LoadConfigFromFile("TestLoadConfigFromFile_GivenNonExistantFilename_ReturnsError")
	assert.Nil(t, conf)
	assert.True(t, os.IsNotExist(err))
}

func TestLoadConfigFromFile_GivenInvalidYAML_ReturnsError(t *testing.T) {
	// Invalid yaml
	const testYAML = `provider: test
migrations:
	- name: Test
	up: test.up.sql
	down: test.down.sql`

	file, err := os.Create("TestLoadConfigFromFile_GivenInvalidYAML_ReturnsError")
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
		return
	}

	file.Write([]byte(testYAML))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestLoadConfigFromFile_GivenInvalidYAML_ReturnsError")
	})

	conf, err := LoadConfigFromFile("TestLoadConfigFromFile_GivenInvalidYAML_ReturnsError")
	assert.Nil(t, conf)
	assert.NotNil(t, err)
}

func TestConfigMapGet_HavingStringValue_ReturnsValue(t *testing.T) {
	conf := ConfigMap{
		"MyString": "Hello World",
	}
	v, ok := conf.String("MyString")
	assert.Equal(t, "Hello World", v)
	assert.True(t, ok)
}

func TestConfigMapGet_GivenNilMap_ReturnsFalse(t *testing.T) {
	var conf ConfigMap = nil
	v, ok := conf.String("MyString")
	assert.Equal(t, "", v)
	assert.False(t, ok)
}

func TestConfigMapGet_ValueDoesNotExist_ReturnsFalse(t *testing.T) {
	conf := make(ConfigMap)
	v, ok := conf.String("MyString")
	assert.Equal(t, "", v)
	assert.False(t, ok)
}

func TestConfigMapGet_ValueIsNotString_ReturnsFalse(t *testing.T) {
	conf := ConfigMap{
		"MyString": 1,
	}
	v, ok := conf.String("MyString")
	assert.Equal(t, "", v)
	assert.False(t, ok)
}
