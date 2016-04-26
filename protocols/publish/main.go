package publish

import (
	"github.com/dappstore/go-dapp"
)

// Protocol represents a configuration of the publish protocol
type Protocol struct {
	store dapp.Store
	kv    dapp.KV
}

// New creates a new dfs protocol
func New(kv dapp.KV, store dapp.Store) *Protocol {
	return &Protocol{kv: kv, store: store}
}
