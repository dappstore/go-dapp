package dfs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// LoadTemp loads `contents` into a temporary path
func (sys *Protocol) LoadTemp(contents dapp.Hash) (string, error) {
	dir, err := ioutil.TempDir("", "dapp-dfs")
	if err != nil {
		return "", errors.Wrap(err, "protocol-dfs: failed to create temp dir")
	}

	err = os.Remove(dir)
	if err != nil {
		return "", errors.Wrap(err, "protocol-dfs: failed to remove temp dir")
	}

	err = sys.store.LoadPath(dir, contents)
	if err != nil {
		return "", errors.Wrap(err, "protocol-dfs: load local dir failed")
	}

	return dir, nil
}

// LoadTempDir loads all hashes into a temp directory
func (sys *Protocol) LoadTempDir(contents map[string]dapp.Hash) (string, error) {

	dir, err := ioutil.TempDir("", "dapp-dfs")
	if err != nil {
		return "", errors.Wrap(err, "protocol-dfs: failed to create temp dir")
	}

	for name, hash := range contents {
		err = sys.store.LoadPath(filepath.Join(dir, name), hash)
		if err != nil {
			return "", errors.Wrap(err, "protocol-dfs: store load failed")
		}
	}

	return dir, nil
}

// MergeAtPath loads `source` into a temporary directory, ensures that `path`
// doesn't exist (and removes it if it does), then adds `contents` at `path`.
func (sys *Protocol) MergeAtPath(
	source dapp.Hash,
	path string,
	contents dapp.Hash,
) (result dapp.Hash, err error) {

	dir, err := sys.LoadTemp(source)
	if err != nil {
		err = errors.Wrap(err, "protocol-dfs: loading temp failed")
		return
	}
	// defer os.RemoveAll(dir)

	dest := filepath.Join(dir, path)
	_, err = os.Stat(dest)
	if err == nil {
		if err != nil {
			err = errors.Wrap(err, "protocol-dfs: failed to remove existing content")
			return
		}
	}

	if !os.IsNotExist(err) {
		err = errors.Wrap(err, "protocol-dfs: failed to stat destination")
		return
	}

	err = os.MkdirAll(filepath.Dir(dest), 0700)
	if err != nil {
		err = errors.Wrap(err, "protocol-dfs: failed to create path for new content")
		return
	}

	err = sys.store.LoadPath(dest, contents)
	if err != nil {
		err = errors.Wrap(err, "protocol-dfs: merging failed")
		return
	}

	result, err = sys.store.StorePath(dir)
	if err != nil {
		err = errors.Wrap(err, "protocol-dfs-merge: store failed")
		return
	}

	return
}

// StoreDir adds `contents` into the store grouped together as a directory
func (sys *Protocol) StoreDir(contents map[string]dapp.Hash) (dapp.Hash, error) {
	s := sys.store

	dir, err := sys.LoadTempDir(contents)
	if err != nil {
		return dapp.Hash{}, errors.Wrap(err, "StoreDir: loading local dir failed")
	}

	defer os.RemoveAll(dir)

	h, err := s.StorePath(dir)
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
		contents[name], err = s.StorePath(path)
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

// StoreString adds `contents` into the store a file and returns its hash
func (sys *Protocol) StoreString(
	contents string,
) (ret dapp.Hash, err error) {
	dir, err := ioutil.TempDir("", "dapp-store-temp")

	if err != nil {
		err = errors.Wrap(err, "protocol-dfs: tempdir failed")
		return
	}

	path := filepath.Join(dir, "contents")
	err = ioutil.WriteFile(path, []byte(contents), 0600)
	if err != nil {
		err = errors.Wrap(err, "protocol-dfs: write contents failed")
		return
	}

	ret, err = sys.store.StorePath(path)
	if err != nil {
		err = errors.Wrap(err, "protocol-dfs: ipfs add failed")
		return
	}

	return
}
