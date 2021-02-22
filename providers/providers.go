package providers

import (
	"fmt"
	"sync"

	"github.com/reecerussell/migrations"
)

var (
	mu   = sync.Mutex{}
	prvs = make(map[string]ConstructorFunc)
)

// ConstructorFunc is used as a build func to build a provider,
// giving it access to a ConfigMap, if needed.
type ConstructorFunc func(migrations.ConfigMap) migrations.Provider

// Add registers a provider by loading a constructur and name
// into memory, which can then be consumer via the Get function.
func Add(name string, p ConstructorFunc) {
	mu.Lock()
	defer mu.Unlock()

	prvs[name] = p
}

// Get retrieves a provider with the given name, by building it
// using the registered ConstructorFunc, providing it with
// the given ConfigMap.
func Get(name string, conf migrations.ConfigMap) migrations.Provider {
	mu.Lock()
	defer mu.Unlock()

	p, ok := prvs[name]
	if !ok {
		panic(fmt.Sprintf("no provider '%s' is registered", name))
	}

	return p(conf)
}
