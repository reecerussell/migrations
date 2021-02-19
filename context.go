package migrations

import (
	"context"
	"time"
)

// Context is an implementation of context.Context.
type Context struct {
	ctx context.Context

	FileContext string
}

// NewContext returns a new instance of Context, implementing context.Context.
func NewContext(ctx context.Context, fileContext string) context.Context {
	return &Context{
		ctx:         ctx,
		FileContext: fileContext,
	}
}

// Deadline calls the base context's Deadline func.
func (mc *Context) Deadline() (time.Time, bool) {
	return mc.ctx.Deadline()
}

// Done calls the base context's Done func.
func (mc *Context) Done() <-chan struct{} {
	return mc.ctx.Done()
}

// Err calls the base context's Err func.
func (mc *Context) Err() error {
	return mc.ctx.Err()
}

// Value calls the base context's Value func.
func (mc *Context) Value(key interface{}) interface{} {
	return mc.ctx.Value(key)
}
