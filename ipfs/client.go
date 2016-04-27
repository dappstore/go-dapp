package ipfs

import (
	"io"
	"log"
	"os"

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
	hash, err := c.StorePath(path)
	if err != nil {
		panic(errors.Wrap(err, "ipfs-hasher: store failed"))
	}

	return hash
}

// LoadPath implements dapp.Store
func (c *Client) LoadPath(dir string, content dapp.Hash) error {
	_, err := os.Stat(dir)
	if err == nil {
		return errors.New("ipfs-load: destination exists")
	}

	if !os.IsNotExist(err) {
		return errors.Wrap(err, "ipfs: stat destination failed")
	}

	ipfsPath := Join(content.Multihash)

	log.Println("get", ipfsPath, dir)
	err = c.shell.Get(ipfsPath, dir)
	if err != nil {
		return errors.Wrap(err, "ipfs: get failed")
	}

	return nil
}

// StorePath implements dapp.Store
func (c *Client) StorePath(path string) (dapp.Hash, error) {
	log.Println("ISSUE:  for some reason, if path is a directory, it gets nested within itself.  seems to be a the go libs issue, since cli works")
	hash, err := c.add(path)
	return dapp.Hash{Multihash: hash}, err
}

func (c *Client) add(path string) (multihash.Multihash, error) {

	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "ipfs-add: path doesn't exist")
	}

	var hashStr string

	if stat.IsDir() {
		hashStr, err = c.shell.AddDir(path)
	} else {
		var f io.ReadCloser
		f, err = os.Open(path)
		if err != nil {
			return nil, errors.Wrap(err, "ipfs-add: couldn't open path")
		}
		defer f.Close()

		hashStr, err = c.shell.Add(f)
	}

	if err != nil {
		return nil, errors.Wrap(err, "ipfs-add: shell command failed")
	}

	hash, err := multihash.FromB58String(hashStr)
	if err != nil {
		return nil, errors.Wrap(err, "ipfs: failed to parse add result")
	}

	return hash, nil
}
