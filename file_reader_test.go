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

	fr := NewFileReader()
	bytes, err := fr.Read(".", "FileReader")
	assert.Equal(t, []byte("Hello World"), bytes)
	assert.NoError(t, err)
}

func TestFileReaderRead_GivenInvalidFilePath_ReturnsIsNotExist(t *testing.T) {
	fr := NewFileReader()
	bytes, err := fr.Read(".", "MissingFileName")
	assert.Nil(t, bytes)
	assert.True(t, os.IsNotExist(err))
}
