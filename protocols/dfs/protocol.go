package dfs

import (
	"os"
	"path/filepath"

	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// LoadTemp loads `contents` into a temporary path
func (sys *Protocol) LoadTemp(contents dapp.Hash) (string, error) {
	s := sys.store

	dir, err := s.NewTempDir()
	if err != nil {
		return "", errors.Wrap(err, "protocol-fs: failed to create temp dir")
	}

	err = s.LoadLocalDir(filepath.Join(dir), contents)
	if err != nil {
		return "", errors.Wrap(err, "protocol-fs: load local dir failed")
	}

	return dir, nil
}

// LoadTempDir loads all hashes into a temp directory
func (sys *Protocol) LoadTempDir(contents map[string]dapp.Hash) (string, error) {
	s := sys.store

	dir, err := s.NewTempDir()
	if err != nil {
		return "", errors.Wrap(err, "LoadMap: failed to create local dir")
	}

	for name, hash := range contents {
		err = s.LoadLocalDir(filepath.Join(dir, name), hash)
		if err != nil {
			return "", errors.Wrap(err, "LoadMap: store load failed")
		}
	}

	return dir, nil
}

// StoreDir adds `contents` into the store grouped together as a directory
func (sys *Protocol) StoreDir(contents map[string]dapp.Hash) (dapp.Hash, error) {
	s := sys.store

	dir, err := sys.LoadTempDir(contents)
	if err != nil {
		return dapp.Hash{}, errors.Wrap(err, "StoreDir: loading local dir failed")
	}

	defer os.RemoveAll(dir)

	h, err := s.StoreLocalDir(dir)
	if err != nil {
		return dapp.Hash{}, errors.Wrap(err, "StoreDir: ipfs add failed")
	}

	return h, nil
}

// StoreLocalPaths adds `contents` into the store as groups together as a directory
func (sys *Protocol) StoreLocalPaths(paths []string) (dapp.Hash, error) {
	s := sys.store
	contents := map[string]dapp.Hash{}

	// Add all paths to store, collecting hashes
	for _, path := range paths {
		var err error
		name := filepath.Base(path)
		contents[name], err = s.StoreLocalDir(path)
		if err != nil {
			return dapp.Hash{},
				errors.Wrap(err, "StoreLocalPaths: failed storing child")
		}
	}

	h, err := sys.StoreDir(contents)
	if err != nil {
		return dapp.Hash{}, errors.Wrap(err, "StoreLocalPaths: StoreDir failed")
	}

	return h, nil
}
