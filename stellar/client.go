package stellar

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
	"github.com/stellar/go-stellar-base/keypair"
)

// Set implements kv.Kv
func (c *Client) Set(identity dapp.Identity, key string, value []byte) (dapp.TX, error) {

	return dapp.TX(""), nil
}

// Get implements kv.Kv
func (c *Client) Get(identity dapp.Identity, key string) ([]byte, error) {
	return nil, nil
}

// ParseIdentity implements dapp.IdentityProvider
func (c *Client) ParseIdentity(str string) (dapp.Identity, error) {
	kp, err := keypair.Parse(str)
	if err != nil {
		return nil, errors.Wrap(err, "parse identity")
	}

	return &Identity{KP: kp}, nil
}

// RandomIdentity implements dapp.IdentityProvider
func (c *Client) RandomIdentity() (dapp.Identity, error) {
	kp, err := keypair.Random()
	if err != nil {
		return nil, errors.Wrap(err, "stellar: create random keypair failed")
	}

	return &Identity{KP: kp}, nil
}
