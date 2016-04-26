// Package claim implements the claim protocol.
//
// As part of applying policies within a process, claims can be made.  these
// claims can be serialized to json and recorded. claims can be trusted.
package claim

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Make records a claim made by the current process.
func (p *Protocol) Make(path string, value interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	// if a value exists at the path, error
	if p.Claims.ExistsP(path) {
		return errors.New("claim already set")
	}

	p.Claims.SetP(value, path)
	return nil
}

// Push pushes a claim on to an array at `path` for current process.
func (p *Protocol) Push(path string, value interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	var err error
	if !p.Claims.ExistsP(path) {
		_, err = p.Claims.ArrayP(path)
		if err != nil {
			return errors.Wrap(err, "claim at path is not an array")
		}
	}

	err = p.Claims.ArrayAppendP(value, path)
	if err != nil {
		return errors.Wrap(err, "claim at path is not an array")
	}

	return nil
}

// WriteFile saves the current process' claims to disk
func (p *Protocol) WriteFile(fs afero.Fs, path string, perm os.FileMode) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	err := afero.WriteFile(fs, path, p.Claims.Bytes(), perm)
	if err != nil {
		return errors.Wrap(err, "claim write failed")
	}

	return nil
}
