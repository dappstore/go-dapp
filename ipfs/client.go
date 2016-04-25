package ipfs

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/dappstore/go-dapp"
	"github.com/jbenet/go-multihash"
	"github.com/pkg/errors"
)

// HashLocalPath implements hash.Hasher
func (c *Client) HashLocalPath(path string) dapp.Hash {
	hash, err := add(path)
	if err != nil {
		panic(err)
	}

	return dapp.Hash{Multihash: hash}
}

// LoadLocalDir implements store.ContentAddressableFS
func (c *Client) LoadLocalDir(content dapp.Hash) (string, error) {
	dir, err := ioutil.TempDir("", "dapp-load")

	// TODO: make more clear
	ipfsPath := fmt.Sprintf("/ipfs/%s", content.Multihash.B58String())
	err = exec.Command("ipfs", "get", "-o", dir, ipfsPath).Run()
	if err != nil {
		return "", errors.Wrap(err, "ipfs: get failed")

	}

	return dir, nil
}

// StoreLocalDir implements store.ContentAddressableFS
func (c *Client) StoreLocalDir(path string) (dapp.Hash, error) {
	hash, err := add(path)
	return dapp.Hash{Multihash: hash}, err
}

// Hash returns the hash of `path` according to ipfs
func add(path string) (ret multihash.Multihash, err error) {
	stdout, err := exec.Command("ipfs", "add", "-r", "-q", path).Output()
	if err != nil {
		err = errors.Wrap(err, "ipfs: add failed")
		return
	}

	hashes := strings.Split(strings.TrimSpace(string(stdout)), "\n")
	lastHash := hashes[len(hashes)-1]

	ret, err = multihash.FromB58String(lastHash)
	if err != nil {
		err = errors.Wrap(err, "ipfs: failed decoding ipfs output")
		return
	}

	return
}
