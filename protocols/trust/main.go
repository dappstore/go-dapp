package trust

import (
	"github.com/dappstore/go-dapp"
)

// Protocol represents a configuration of this protocol
type Protocol struct {
	trustSets map[string]Set
}

// Set is a set of identities that are trusted in some way
type Set []dapp.Identity
