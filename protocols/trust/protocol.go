package trust

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// AddTrust records trust in `id` for `role`
func (sys *Protocol) AddTrust(id dapp.Identity, role string) error {
	return nil
}

// RemoveTrust removes trust
func (sys *Protocol) RemoveTrust(id dapp.Identity, role string) error {
	set := sys.trustSets[role]
	toRemove := -1

	for i, trusted := range set {
		if trusted.Equals(id) {
			toRemove = i
			break
		}
	}

	if toRemove == -1 {
		return errors.New("protocol-trust: id is already not trusted")
	}

	var next = append(set[:toRemove], set[toRemove+1:]...)
	sys.trustSets[role] = next

	return nil
}

// IsTrusted returns true if `id` is trusted in `role` according this systems
// local trust set.
func (sys *Protocol) IsTrusted(id dapp.Identity, role string) (bool, error) {
	return false, nil
}
