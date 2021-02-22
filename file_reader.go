package migrations

import (
	"io/ioutil"
	"path"
)

// FileReader is a high level interface used to read files
// in a specific directory. Used to read migration files.
type FileReader interface {
	// Read returns the file's content from the given directory.
	Read(directory, filename string) ([]byte, error)
}

// fileReader is an implementation of FileReader.
type fileReader struct{}

// NewFileReader returns a new instance of FileReader.
func NewFileReader() FileReader {
	return &fileReader{}
}

func (*fileReader) Read(directory, filename string) ([]byte, error) {
	filePath := path.Join(directory, filename)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
