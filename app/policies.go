package app

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// IdentityProvider is a policy that registers an identity system.
type IdentityProvider struct {
	dapp.IdentityProvider
}

var _ Policy = &IdentityProvider{}

// ApplyDappPolicy implements `Policy`
func (p *IdentityProvider) ApplyDappPolicy(app *App) error {
	if app.Providers.IdentityProvider != nil {
		return errors.New("policy: cannot overwrite identity system")
	}

	if p.IdentityProvider == nil {
		return errors.New("policy: cannot apply nil identity system")
	}

	app.Providers.IdentityProvider = p.IdentityProvider
	return nil
}

// KV is a policy that registers a decentralized key value store when applied.
type KV struct {
	dapp.KV
}

var _ Policy = &KV{}

// ApplyDappPolicy implements `Policy`
func (p *KV) ApplyDappPolicy(app *App) error {
	if app.Providers.KV != nil {
		return errors.New("policy: cannot overwrite kv system")
	}

	if p.KV == nil {
		return errors.New("policy: cannot apply nil kv system")
	}

	app.Providers.KV = p.KV
	return nil
}

// RunVerification represents the dapp policy that actually runs the process
// verification protocol.
type RunVerification struct{}

var _ Policy = &RunVerification{}

// ApplyDappPolicy applies `p` to `app`
func (p *RunVerification) ApplyDappPolicy(app *App) error {
	return nil
}

// Store is a policy that registers a content addressable store when applied
type Store struct {
	dapp.Store
}

var _ Policy = &Store{}

// ApplyDappPolicy implements `Policy`
func (p *Store) ApplyDappPolicy(app *App) error {
	if app.Providers.Store != nil {
		return errors.New("policy: cannot overwrite store system")
	}

	if p.Store == nil {
		return errors.New("policy: cannot apply nil store system")
	}

	app.Providers.Store = p.Store
	return nil
}

// VerifySelf is a policy that causes the binary to verify itself as an
// installation of the application published by `Publisher`, according to the
// dapp publisher protocol.
type VerifySelf struct {
	Publisher string
}

var _ Policy = &VerifySelf{}

// ApplyDappPolicy applies `p` to `app`
func (p *VerifySelf) ApplyDappPolicy(app *App) error {
	return nil
}
