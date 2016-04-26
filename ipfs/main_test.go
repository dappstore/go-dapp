package ipfs_test

import (
	"github.com/dappstore/go-dapp"
	"github.com/dappstore/go-dapp/ipfs"
	"github.com/dappstore/go-dapp/protocols/claim"
	"github.com/dappstore/go-dapp/protocols/hash"
)

var _ hash.Hasher = ipfs.DefaultClient
var _ claim.MakesClaims = ipfs.DefaultClient
var _ dapp.Store = ipfs.DefaultClient
