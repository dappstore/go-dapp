package ipfs

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dappstore/go-dapp"
	"github.com/jbenet/go-multihash"
	"github.com/pkg/errors"
)

// ClaimIdentity is the dapp identity for this package
const ClaimIdentity = "GAPFBAPSCKBJH6HRDXFIMOI367L3T3SMJ5I2EBW4BVVXSCL2WYNTJ5WL"

// Add ensures `path` is in ipfs
func (c *Client) Add(path string) (multihash.Multihash, error) {
	hashStr, err := c.shell.AddDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "ipfs: add failed")
	}

	hash, err := multihash.FromB58String(hashStr)
	if err != nil {
		return nil, errors.Wrap(err, "ipfs: failed to parse add result")
	}

	return hash, nil
}

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

// Get loads the ipfs data in the directory `dir` underneath `base` into the
// local directory at `local`.
func (c *Client) Get(base multihash.Multihash, dir string, local string) error {

	// TODO: make more clear
	ipfsPath := Join(base, dir)
	if dir == "" {
		ipfsPath = Join(base)
	}

	err := c.shell.Get(ipfsPath, local)
	if err != nil {
		return errors.Wrap(err, "ipfs: get failed")
	}

	return nil
}

// Exists checks to see if `base` has a child named `child` in ipfs
func (c *Client) Exists(base multihash.Multihash, child string) (bool, error) {
	ipfsPath := Join(base, child)
	_, err := c.shell.List(ipfsPath)

	if err != nil {
		return false, errors.Wrap(err, "ipfs: list failed")
	}

	return true, nil
}

// HashLocalPath implements hash.Hasher
func (c *Client) HashLocalPath(path string) dapp.Hash {
	hash, err := c.Add(path)
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
	ipfsPath := Join(content.Multihash)

	log.Println("get", ipfsPath, dir)
	// err = c.shell.Get(ipfsPath, dir)
	// if err != nil {
	// 	return errors.Wrap(err, "ipfs: get failed")
	// }

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
	hash, err := c.Add(path)
	return dapp.Hash{Multihash: hash}, err
}
