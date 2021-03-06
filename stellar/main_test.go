package stellar_test

import (
	"github.com/dappstore/go-dapp"
	"github.com/dappstore/go-dapp/protocols/claim"
	"github.com/dappstore/go-dapp/stellar"
)

var _ dapp.Identity = &stellar.Identity{}
var _ claim.MakesClaims = stellar.DefaultClient
var _ dapp.KV = stellar.DefaultClient
var _ dapp.IdentityProvider = stellar.DefaultClient
