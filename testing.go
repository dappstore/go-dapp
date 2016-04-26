package dapp

import (
	"github.com/pkg/errors"
)

// ensure our mocks implement our interfaces
var _ Identity = &MockIdentity{}

//MockIdentity is a mock identity.  use it in your tests that are dependent upon
//this package.
type MockIdentity struct {
	PK string
}

// Equals implements dapp.Identity
func (i *MockIdentity) Equals(other Identity) bool {
	oid, ok := other.(*MockIdentity)
	if !ok {
		return false
	}

	return i.PK == oid.PK
}

// PublicKey implement `Identity`
func (i *MockIdentity) PublicKey() string {
	return i.PK
}

// Verify implement `Identity`
func (i *MockIdentity) Verify(input []byte, signature []byte) error {
	return errors.New("mock identity cannot verify signatures")
}

// Sign implement `Identity`
func (i *MockIdentity) Sign(input []byte) ([]byte, error) {
	return nil, errors.New("mock identity cannot sign messages")
}
