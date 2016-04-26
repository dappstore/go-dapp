package ipfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/dappstore/go-dapp"
	"github.com/jbenet/go-multihash"
	"github.com/pkg/errors"
)

// ClaimIdentity is the dapp identity for this package
const ClaimIdentity = "GAPFBAPSCKBJH6HRDXFIMOI367L3T3SMJ5I2EBW4BVVXSCL2WYNTJ5WL"

// ClaimerName implements `MakesClaims`
func (c *Client) ClaimerName() string {
	return "ipfs"
}

// ClaimerIdentity implements `MakesClaims`
func (c *Client) ClaimerIdentity() string {
	return ClaimIdentity
}

// ClaimerClaims implements `MakesClaims`
func (c *Client) ClaimerClaims() string { return "" }

// HashLocalPath implements hash.Hasher
func (c *Client) HashLocalPath(path string) dapp.Hash {
	hash, err := add(path)
	if err != nil {
		panic(err)
	}

	return dapp.Hash{Multihash: hash}
}

// LoadLocalDir implements dapp.Store
func (c *Client) LoadLocalDir(dir string, content dapp.Hash) error {
	stat, err := os.Stat(dir)

	// if the error is not "does not exists" and is populated, error out
	if !os.IsNotExist(err) && err != nil {
		return errors.Wrap(err, "ipfs: stat destination failed")
	}

	// if it exists, ensure it is a dir
	if !os.IsNotExist(err) && !stat.IsDir() {
		return errors.New("ipfs: destination is not directory")
	}

	// TODO: make more clear
	ipfsPath := fmt.Sprintf("/ipfs/%s", content.Multihash.B58String())
	err = exec.Command("ipfs", "get", "-o", dir, ipfsPath).Run()
	if err != nil {
		return errors.Wrap(err, "ipfs: get failed")

	}

	return nil
}

// NewTempDir implements dapp.Store
func (c *Client) NewTempDir() (string, error) {
	dir, err := ioutil.TempDir("", "dapp-store-temp")

	if err != nil {
		return "", errors.Wrap(err, "ipfs: tempdir failed")

	}

	return dir, nil
}

// StoreLocalDir implements dapp.Store
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
