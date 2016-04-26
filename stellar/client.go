package stellar

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
	"github.com/stellar/go-stellar-base/build"
	"github.com/stellar/go-stellar-base/keypair"
)

// ClaimIdentity is the dapp identity for this package
const ClaimIdentity = "GBHE5TFF6PLBUCO4BTBX6L5AX6CPG6ACHRMA6UJ7VTKCT674UVNV5LA2"

// ClaimerName implements `MakesClaims`
func (c *Client) ClaimerName() string { return "stellar" }

// ClaimerIdentity implements `MakesClaims`
func (c *Client) ClaimerIdentity() string {
	return ClaimIdentity
}

// ClaimerClaims implements `MakesClaims`
func (c *Client) ClaimerClaims() string { return "" }

// Set implements kv.Kv
func (c *Client) Set(identity dapp.Identity, key string, value []byte) (dapp.TX, error) {

	sid := identity.(*Identity)

	full, ok := sid.KP.(*keypair.Full)
	if !ok {
		return dapp.TX(""), errors.New("stellar: don't know secret key for identity")
	}

	tx := build.Transaction(
		build.SourceAccount{AddressOrSeed: sid.PublicKey()},
		build.AutoSequence{SequenceProvider: c.Client},
		build.SetData("dapp:publications", value),
	)

	txe := tx.Sign(full.Seed())

	xdrs, err := txe.Base64()
	if err != nil {
		return dapp.TX(""), errors.Wrap(err, "stellar: failed to craft transaction")
	}

	result, err := c.Client.SubmitTransaction(xdrs)
	if err != nil {
		return dapp.TX(""), errors.Wrap(err, "stellar: transaction failed")
	}

	return dapp.TX(result.Hash), nil
}

// Get implements kv.Kv
func (c *Client) Get(identity dapp.Identity, key string) ([]byte, error) {
	sid := identity.(*Identity)
	data, err := LoadAccountData(c.Client, sid.Address())
	if err != nil {
		return nil, errors.Wrap(err, "stellar: load account failed")
	}

	return data[key], nil
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

// AnnounceIdentity implements dapp.IdentityProvider
func (c *Client) AnnounceIdentity(id dapp.Identity) (dapp.TX, error) {
	sid := id.(*Identity)
	txHash, err := FundAccount(c.Client, sid.Address())
	if err != nil {
		return dapp.TX(""), errors.Wrap(err, "stellar: funding account failed")
	}

	return dapp.TX(txHash), nil

}

// IsIdentityAnnounced implements dapp.IdentityProvider
func (c *Client) IsIdentityAnnounced(id dapp.Identity) (bool, error) {
	sid := id.(*Identity)
	return AccountExists(c.Client, sid.Address())
}
