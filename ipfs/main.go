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

// Exists checks to see if `base` has a child named `child` in ipfs
func Exists(base multihash.Multihash, child string) (bool, error) {
	return DefaultClient.Exists(base, child)
}

// Join produces a new ipfs path from `base` and `dirs`
func Join(base multihash.Multihash, dirs ...string) string {
	return strings.Join(
		append([]string{fmt.Sprintf("/ipfs/%s", base.B58String())}, dirs...),
		"/",
	)
}
