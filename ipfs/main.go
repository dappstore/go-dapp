package ipfs

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/jbenet/go-multihash"
	"github.com/pkg/errors"
)

// DefaultClient is the default client
var DefaultClient = &Client{}

// Client uses the ipfs cli app
type Client struct{}

// Add ensures `path` is in ipfs
func Add(path string) (multihash.Multihash, error) {
	return Hash(path)
}

// Exists checks to see if `base` has a child named `child` in ipfs
func Exists(base multihash.Multihash, child string) (bool, error) {
	ipfsPath := Join(base, child)
	_, err := exec.Command("ipfs", "ls", ipfsPath).Output()

	if eerr, ok := err.(*exec.ExitError); ok {
		exit := eerr.Sys().(syscall.WaitStatus)

		if exit.ExitStatus() == 1 {
			return false, nil
		}
	}

	if err != nil {
		return false, errors.Wrap(err, "ipfs ls failed")
	}

	return true, nil
}

// Get loads the ipfs data in the directory `dir` underneath `base` into the
// local directory at `local`.
func Get(base multihash.Multihash, dir string, local string) error {

	// TODO: make more clear
	ipfsPath := Join(base, dir)
	if dir == "" {
		ipfsPath = Join(base)
	}
	err := exec.Command("ipfs", "get", "-o", local, ipfsPath).Run()
	if err != nil {
		return errors.Wrap(err, "ipfs get failed")

	}

	return nil
}

// Hash returns the hash of `path` according to ipfs
func Hash(path string) (ret multihash.Multihash, err error) {
	stdout, err := exec.Command("ipfs", "add", "-r", "-q", path).Output()
	if err != nil {
		err = errors.Wrap(err, "ipfs add failed")
		return
	}

	hashes := strings.Split(strings.TrimSpace(string(stdout)), "\n")
	lastHash := hashes[len(hashes)-1]

	ret, err = multihash.FromB58String(lastHash)
	if err != nil {
		err = errors.Wrap(err, "failed decoding ipfs output")
		return
	}

	return
}

// Join produces a new ipfs path from `base` and `dirs`
func Join(base multihash.Multihash, dirs ...string) string {
	return fmt.Sprintf("/ipfs/%s/%s",
		base.B58String(),
		strings.Join(dirs, "/"),
	)
}
