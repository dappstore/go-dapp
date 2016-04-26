package publish

import (
	"github.com/dappstore/go-dapp"
	"github.com/dappstore/go-dapp/protocols/claim"
)

// Protocol represents a configuration of the publish protocol
type Protocol struct {
	claims *claim.Protocol
	store  dapp.Store
	kv     dapp.KV
}
