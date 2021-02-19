package migrations

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	testCtx := context.Background()
	testFileContext := "."

	ctx := NewContext(testCtx, testFileContext)
	assert.Equal(t, testCtx, ctx.(*Context).ctx)
	assert.Equal(t, testFileContext, ctx.(*Context).FileContext)
}

func TestMigrationContextDeadline_GivenCtxWithDeadline_ReturnsDeadline(t *testing.T) {
	testDeadline := time.Now()
	testCtx, cancel := context.WithDeadline(context.Background(), testDeadline)
	defer cancel()

	ctx := &Context{ctx: testCtx}
	deadline, ok := ctx.Deadline()
	assert.Equal(t, testDeadline, deadline)
	assert.True(t, ok)
}

func TestMigrationContextDone_GivenCtxWithCancel_ReturnsCancelChan(t *testing.T) {
	testCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx := &Context{ctx: testCtx}
	assert.Equal(t, testCtx.Done(), ctx.Done())
}

func TestMigrationContextErr_GivenEmptyContext_ReturnsNoError(t *testing.T) {
	ctx := &Context{ctx: context.Background()}
	assert.NoError(t, ctx.Err())
}

func TestMigrationContextValue_GivenValidKey_ReturnsValue(t *testing.T) {
	testKey := Migration{}
	testValue := "MyMigration"

	testCtx := context.WithValue(context.Background(), testKey, testValue)
	ctx := &Context{ctx: testCtx}

	value := ctx.Value(testKey)
	assert.Equal(t, value, testValue)
}
