package migrations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationUp_WithValidFile_ReturnsFileContent(t *testing.T) {
	const testData = "Migration Data"

	file, err := os.Create("TestMigrationUp_WithValidFile_ReturnsFileContent")
	if err != nil {
		t.Errorf("failed to create test file: %v", err)
		return
	}

	file.Write([]byte(testData))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestMigrationUp_WithValidFile_ReturnsFileContent")
	})

	ctx := &Context{FileContext: "."}
	m := &Migration{UpFile: "TestMigrationUp_WithValidFile_ReturnsFileContent"}

	content, err := m.Up(ctx)
	assert.NoError(t, err)
	assert.Equal(t, testData, content)
}

func TestMigrationUp_WithNonExistantFile_ReturnsError(t *testing.T) {
	ctx := &Context{FileContext: "."}
	m := &Migration{UpFile: "TestMigrationUp_WithNonExistantFile_ReturnsError"}

	content, err := m.Up(ctx)
	assert.Equal(t, "", content)
	assert.True(t, os.IsNotExist(err))
}

func TestMigrationUp_WithNonExistantFileContext_ReturnsError(t *testing.T) {
	ctx := &Context{FileContext: "./test_migration_path"}
	m := &Migration{UpFile: "TestMigrationUp_WithNonExistantFileContext_ReturnsError"}

	content, err := m.Up(ctx)
	assert.Equal(t, "", content)
	assert.True(t, os.IsNotExist(err))
}

func TestMigrationDown_WithValidFile_ReturnsFileContent(t *testing.T) {
	const testData = "Migration Data"

	file, err := os.Create("TestMigrationDown_WithValidFile_ReturnsFileContent")
	if err != nil {
		t.Errorf("failed to create test file: %v", err)
		return
	}

	file.Write([]byte(testData))
	file.Close()

	t.Cleanup(func() {
		os.Remove("TestMigrationDown_WithValidFile_ReturnsFileContent")
	})

	ctx := &Context{FileContext: "."}
	m := &Migration{DownFile: "TestMigrationDown_WithValidFile_ReturnsFileContent"}

	content, err := m.Down(ctx)
	assert.NoError(t, err)
	assert.Equal(t, testData, content)
}

func TestMigrationDown_WithNonExistantFile_ReturnsError(t *testing.T) {
	ctx := &Context{FileContext: "."}
	m := &Migration{DownFile: "TestMigrationDown_WithNonExistantFile_ReturnsError"}

	content, err := m.Down(ctx)
	assert.Equal(t, "", content)
	assert.True(t, os.IsNotExist(err))
}

func TestMigrationDown_WithNonExistantFileContext_ReturnsError(t *testing.T) {
	ctx := &Context{FileContext: "./test_migration_path"}
	m := &Migration{DownFile: "TestMigrationDown_WithNonExistantFileContext_ReturnsError"}

	content, err := m.Down(ctx)
	assert.Equal(t, "", content)
	assert.True(t, os.IsNotExist(err))
}
