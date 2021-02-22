package migrations

import (
	"io/ioutil"
	"path"
)

// FileReader is a high level interface used to read files
// in a specific directory. Used to read migration files.
type FileReader interface {
	Read(filename string) (string, error)
}

// fileReader is an implementation of FileReader.
type fileReader struct {
	fileContext string
}

// NewFileReader returns a new instance of FileReader.
func NewFileReader(fileContext string) FileReader {
	return &fileReader{fileContext: fileContext}
}

func (fr *fileReader) Read(filename string) (string, error) {
	filePath := path.Join(fr.fileContext, filename)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
