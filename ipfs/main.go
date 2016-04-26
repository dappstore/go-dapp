package ipfs

import (
	"fmt"
	"strings"

	iapi "github.com/ipfs/go-ipfs-api"
	"github.com/jbenet/go-multihash"
)

// DefaultClient is the default client
var DefaultClient = New()

// Client uses the ipfs cli app
type Client struct {
	shell *iapi.Shell
}

// New creates a new ipfs client
func New() *Client {
	return &Client{iapi.NewShell("localhost:5001")}
}

// Add ensures `path` is in ipfs
func Add(path string) (multihash.Multihash, error) {
	return DefaultClient.Add(path)
}

// Exists checks to see if `base` has a child named `child` in ipfs
func Exists(base multihash.Multihash, child string) (bool, error) {
	return DefaultClient.Exists(base, child)
}

// Get loads the ipfs data in the directory `dir` underneath `base` into the
// local directory at `local`.
func Get(base multihash.Multihash, dir string, local string) error {
	return DefaultClient.Get(base, dir, local)
}

// Join produces a new ipfs path from `base` and `dirs`
func Join(base multihash.Multihash, dirs ...string) string {
	return strings.Join(
		append([]string{fmt.Sprintf("/ipfs/%s", base.B58String())}, dirs...),
		"/",
	)
}
