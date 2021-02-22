package migrations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileReaderRead(t *testing.T) {
	file, err := os.Create("FileReader")
	if err != nil {
		panic(err)
	}

	file.Write([]byte("Hello World"))
	file.Close()

	t.Cleanup(func() {
		os.Remove("FileReader")
	})

	fr := NewFileReader(".")
	content, err := fr.Read("FileReader")
	assert.Equal(t, "Hello World", content)
	assert.NoError(t, err)
}

func TestFileReaderRead_GivenInvalidFilePath_ReturnsIsNotExist(t *testing.T) {
	fr := NewFileReader(".")
	content, err := fr.Read("MissingFileName")
	assert.Equal(t, "", content)
	assert.True(t, os.IsNotExist(err))
}
