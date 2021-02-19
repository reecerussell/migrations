package migrations

import (
	"context"
	"io/ioutil"
	"path"
	"time"
)

// Migration represents a migration.
type Migration struct {
	ID          int
	Name        string `yaml:"name"`
	DateApplied time.Time
	UpFile      string `yaml:"up"`
	DownFile    string `yaml:"down"`
}

// Up returns the content of the migrations UpFile, using the
// given context's file context.
func (m *Migration) Up(ctx context.Context) (string, error) {
	fileContext := ctx.(*Context).FileContext
	filePath := path.Join(fileContext, m.UpFile)

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Down returns the content of the migrations DownFile, using the
// given context's file context.
func (m *Migration) Down(ctx context.Context) (string, error) {
	fileContext := ctx.(*Context).FileContext
	filePath := path.Join(fileContext, m.DownFile)

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
