package dfs

import (
	"github.com/dappstore/go-dapp"
)

// Default represents the default instantiation of this protocol
var Default *Protocol

// Protocol represents a configuration of the dfs protocol
type Protocol struct {
	store dapp.Store
}

// New creates a new dfs protocol
func New(store dapp.Store) *Protocol {
	return &Protocol{store: store}
}
