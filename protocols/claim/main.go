// Package claim implements the claim protocol.
//
// As part of applying policies within a process, claims can be made.  these
// claims can be serialized to json and recorded. claims can be trusted by
// identities.
package claim

import (
	"os"
	"sync"

	"github.com/Jeffail/gabs"
	"github.com/spf13/afero"
)

// Default represents the default instantiation of this protocol
var Default *Protocol

// Protocol represents a single set of claims made by the running application.
type Protocol struct {
	Claims *gabs.Container

	lock sync.Mutex
}

// New creates a new claim protocol
func New() *Protocol {
	return &Protocol{Claims: gabs.New()}
}

// Make makes a claim on the default claim protocol
func Make(path string, value interface{}) error {
	return Default.Make(path, value)
}

// Push pushes a claim on the default claim protocol
func Push(path string, value interface{}) error {
	return Default.Make(path, value)
}

// WriteFile writes the claims made on the default claim protocol to disk.
func WriteFile(fs afero.Fs, path string, perm os.FileMode) error {
	return Default.WriteFile(fs, path, perm)
}
