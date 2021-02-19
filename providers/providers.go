package providers

import (
	"fmt"
	"sync"

	"github.com/reecerussell/migrations"
)

var (
	mu   = sync.Mutex{}
	prvs = make(map[string]migrations.Provider)
)

// Add loads a migration provider into memory, allowing it to then
// be used to provide migrations.
func Add(name string, p migrations.Provider) {
	mu.Lock()
	defer mu.Unlock()

	prvs[name] = p
}

// Get retrieves a provider with the given name.
func Get(name string) migrations.Provider {
	mu.Lock()
	defer mu.Unlock()

	p, ok := prvs[name]
	if !ok {
		panic(fmt.Sprintf("no provider '%s' is registered", name))
	}

	return p
}
