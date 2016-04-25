package dapp

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dappstore/agree"
	"github.com/dappstore/dapp/dapp/ipfs"
	"github.com/dappstore/dapp/dapp/stellar"
	"github.com/pkg/errors"
)

// ApplyPolicy applies `p` t `a`
func (a *App) ApplyPolicy(p Policy) error {
	err := p.ApplyDappPolicy(a)
	if err != nil {
		return errors.Wrap(err, "failed applying policy")
	}
	a.policies = append(a.policies, p)
	return nil
}

// ApplyPolicies applies all `policies` onto `a`.
func (a *App) ApplyPolicies(policies ...Policy) error {
	for _, p := range policies {
		a.do(
			func() error { return a.ApplyPolicy(p) },
		)
	}
	return nil
}

// CurrentUser returns the current user's identity
func (a *App) CurrentUser() Identity {
	return loginSessions[a.ID]
}

// InitializePolicies applies `policies` t `a`
func (a *App) InitializePolicies(policies []Policy) error {

	// Apply before policies
	a.ApplyPolicies(
		&VerifySelf{Publisher: a.ID},
		&AgreementPolicy{agree.RequireOracle{}},
	)

	// Apply app policies
	a.ApplyPolicies(policies...)

	// Apply default policies
	if len(a.agreements.Oracles) == 0 {
		a.ApplyPolicies(
			&AgreementOracle{stellar.HorizonOracle(defaultHorizon)},
		)
	}

	// Apply after policies
	a.ApplyPolicies(&RunVerification{})

	return a.Err
}

// Fund funds `user`
func (a *App) Fund(id Identity) error {
	return Fund(id)
}

// LoadMap loads all hashes into a temp directory
func (a *App) LoadMap(contents map[string]Hash) (string, error) {
	dir, err := ioutil.TempDir("", "dapp-load-map")
	if err != nil {
		return "", errors.Wrap(err, "LoadMap: create temp dir failed")
	}

	for name, hash := range contents {
		err = ipfs.Get(hash.Multihash, "", filepath.Join(dir, name))
		if err != nil {
			return "", errors.Wrap(err, "LoadMap: ipfs get failed")
		}
	}

	return dir, nil
}

// Login logs `user` into `a`
func (a *App) Login(user Identity) {
	Login(a.ID, user)
}

// PublishHash publishes `hash` using the user current user.
func (a *App) PublishHash(hash Hash) (string, error) {
	return PublishHash(hash, a.CurrentUser())
}

// PublishMap publishes `contents` as a directory using the current user.
func (a *App) PublishMap(contents map[string]Hash) (string, Hash, error) {
	hash, err := a.StoreMap(contents)
	if err != nil {
		return "", Hash{}, errors.Wrap(err, "PublishMap: store failed")
	}

	tx, err := a.PublishHash(hash)
	if err != nil {
		return "", Hash{}, errors.Wrap(err, "PublishMap: publish failed")
	}

	return tx, hash, nil
}

// StoreMap adds `contents` into ipfs as a directory
func (a *App) StoreMap(contents map[string]Hash) (Hash, error) {
	dir, err := a.LoadMap(contents)
	if err != nil {
		return Hash{}, errors.Wrap(err, "StoreMap: loading local dir failed")
	}
	defer os.RemoveAll(dir)

	h, err := ipfs.Add(dir)
	if err != nil {
		return Hash{}, errors.Wrap(err, "StoreMap: ipfs add failed")
	}

	return Hash{h}, nil
}

// StorePath adds `path` into ipfs, returning it's hash
func (a *App) StorePath(path string) (Hash, error) {
	h, err := ipfs.Add(path)
	if err != nil {
		return Hash{}, errors.Wrap(err, "store-path: ipfs add failed")
	}

	return Hash{h}, nil
}

// do is a helper to only perform actions while an app remains "un-errored".
func (a *App) do(fns ...func() error) error {
	if a.Err != nil {
		return a.Err
	}

	for _, fn := range fns {
		a.Err = fn()
		if a.Err != nil {
			return a.Err
		}
	}
	return nil
}

func (a *App) verifyPublished() error {
	// load the "dapp-manifest" data field for the account at address from all the
	// horizon servers.  If they disagree, fatally error out (in the future,
	// perhaps retry).
	hash, err := ManifestHash(a.ID, a.verificationServers...)
	if err != nil {
		errors.Print(err)
		os.Exit(1)
	}

	_ = hash
	log.Println(hash.B58String())
	return nil
}
