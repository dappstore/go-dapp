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

var Default *Protocol

type Protocol struct {
	Claims *gabs.Container

	lock sync.Mutex
}

func New() *Protocol {
	return &Protocol{Claims: gabs.New()}
}

func Make(path string, value interface{}) error {
	return Default.Make(path, value)
}

func Push(path string, value interface{}) error {
	return Default.Make(path, value)
}

func WriteFile(fs afero.Fs, path string, perm os.FileMode) error {
	return Default.WriteFile(fs, path, perm)
}
