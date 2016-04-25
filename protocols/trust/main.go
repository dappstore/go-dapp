package trust

import (
	"github.com/dappstore/go-dapp"
)

// System represents a configuration of this protocol
type System struct {
	App *dapp.App

	trustSets map[string]Set
}

// Set is a set of identities that are trusted in some way
type Set []dapp.Identity
