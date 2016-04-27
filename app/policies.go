package app

import (
	"github.com/dappstore/go-dapp"
	"github.com/dappstore/go-dapp/ipfs"
	"github.com/dappstore/go-dapp/protocols/claim"
	"github.com/dappstore/go-dapp/stellar"
	"github.com/pkg/errors"
)

// Defaults is a policy that sets the core providers to the dapp system to the
// defaults, namely using stellar for id and kv providers, and ipfs for the
// store provider.
var Defaults = NewPolicy("defaults",
	&Store{Store: ipfs.DefaultClient},
	&KV{KV: stellar.DefaultClient},
	&IdentityProvider{IdentityProvider: stellar.DefaultClient},
)

// Description is a policy that claims the application's description as `desc`
func Description(desc string) Policy {
	return &fnPolicy{
		"set-description",
		func(app *App) error {
			err := claim.Make("dapp.description", desc)
			if err != nil {
				return errors.Wrap(err, "set-developer: failed to claim dapp.description")
			}

			return nil
		},
	}
}

// Developer is a policy that claims the developer identity is `id`.
func Developer(id string) Policy {
	return &fnPolicy{
		"set-developer",
		func(app *App) error {

			ids := app.Providers
			did, err := ids.ParseIdentity(id)
			if app.Providers.IdentityProvider != nil {
				return errors.Wrap(err, "set-developer: failed to parse id")
			}

			err = claim.Make("dapp.developer", did)
			if err != nil {
				return errors.Wrap(err, "set-developer: failed to claim dapp.developer")
			}

			return nil

		},
	}
}

// Name is a policy that claims
func Name(name string) Policy {
	return &fnPolicy{
		"set-name",
		func(app *App) error {
			err := claim.Make("dapp.name", name)
			if err != nil {
				return errors.Wrap(err, "set-name: failed to claim dapp.name")
			}

			return nil
		},
	}
}

// IdentityProvider is a policy that registers an identity system.
type IdentityProvider struct {
	dapp.IdentityProvider
}

// ApplyDappPolicy implements `Policy`
func (p *IdentityProvider) ApplyDappPolicy(app *App) error {
	if app.Providers.IdentityProvider != nil {
		return errors.New("policy: cannot overwrite identity system")
	}

	if p.IdentityProvider == nil {
		return errors.New("policy: cannot apply nil identity system")
	}

	app.Providers.IdentityProvider = p.IdentityProvider
	return addClaimer(p.IdentityProvider)
}

// KV is a policy that registers a decentralized key value store when applied.
type KV struct {
	dapp.KV
}

// ApplyDappPolicy implements `Policy`
func (p *KV) ApplyDappPolicy(app *App) error {
	if app.Providers.KV != nil {
		return errors.New("policy: cannot overwrite kv system")
	}

	if p.KV == nil {
		return errors.New("policy: cannot apply nil kv system")
	}

	app.Providers.KV = p.KV
	return addClaimer(p.KV)
}

// RunVerification represents the dapp policy that actually runs the process
// verification protocol.
type RunVerification struct{}

// ApplyDappPolicy applies `p` to `app`
func (p *RunVerification) ApplyDappPolicy(app *App) error {
	return nil
}

// Store is a policy that registers a content addressable store when applied
type Store struct {
	dapp.Store
}

// ApplyDappPolicy implements `Policy`
func (p *Store) ApplyDappPolicy(app *App) error {
	if app.Providers.Store != nil {
		return errors.New("policy: cannot overwrite store system")
	}

	if p.Store == nil {
		return errors.New("policy: cannot apply nil store system")
	}

	// TODO
	// err := claim.Make("dapp.providers.store", p.Store.Identity())
	// if err != nil {
	// 	return errors.Wrap(err, "store-policy: could not main claim")
	// }

	app.Providers.Store = p.Store
	return addClaimer(p.Store)
}

// VerifySelf is a policy that causes the binary to verify itself as an
// installation of the application published by `Publisher`, according to the
// dapp publisher protocol.
type VerifySelf struct {
	Publisher string
}

// ApplyDappPolicy applies `p` to `app`
func (p *VerifySelf) ApplyDappPolicy(app *App) error {
	claim.Push("dapp.vulnerabilities", "policy-verify-self: not implemented")
	return nil
}

func addClaimer(c interface{}) error {
	claimer, ok := c.(claim.MakesClaims)
	if !ok {
		return nil
	}

	err := claim.AddClaimer(claimer)
	if err != nil {
		return errors.Wrap(err, "store-policy: failed to add claimer")
	}

	return nil
}
