package stellar

import (
	"github.com/dappstore/go-dapp"
)

// Equals implements dapp.Identity
func (i *Identity) Equals(other dapp.Identity) bool {
	oid, ok := other.(*Identity)
	if !ok {
		return false
	}

	return i.Address() == oid.Address()
}

// PublicKey implements dapp.Identity
func (i *Identity) PublicKey() string {
	return i.KP.Address()
}

// String implements Stringer
func (i *Identity) String() string {
	return i.KP.Address()
}
