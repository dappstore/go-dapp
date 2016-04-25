package dapp

import (
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
		err := a.ApplyPolicy(p)
		if err != nil {
			return errors.Wrap(err, "dapp: failed-policy")
		}
	}

	return nil
}

// CurrentUser returns the current user's identity
func (a *App) CurrentUser() Identity {
	return loginSessions[a.ID]
}

// // InitializePolicies applies `policies` t `a`
// func (a *App) InitializePolicies(policies []Policy) error {
//
// 	// Apply before policies
// 	a.ApplyPolicies(
// 		&VerifySelf{Publisher: a.ID},
// 		&AgreementPolicy{agree.RequireOracle{}},
// 	)
//
// 	// Apply app policies
// 	a.ApplyPolicies(policies...)
//
// 	// Apply default policies
// 	if len(a.agreements.Oracles) == 0 {
// 		a.ApplyPolicies(
// 			&AgreementOracle{stellar.HorizonOracle(defaultHorizon)},
// 		)
// 	}
//
// 	// Apply after policies
// 	a.ApplyPolicies(&RunVerification{})
//
// 	return a.Err
// }

// Fund funds `user`
// func (a *App) Fund(id Identity) error {
// 	return Fund(id)
// }

// KV returns a kv
func (a *App) KV() KV {
	return a.Providers.KV
}

// Login logs `user` into `a`
func (a *App) Login(user Identity) {
	Login(a.ID, user)
}

// Store returns a store
func (a *App) Store() Store {
	return a.Providers.Store
}

// PublishHash publishes `hash` using the user current user.
// func (a *App) PublishHash(hash Hash) (string, error) {
// 	return PublishHash(hash, a.CurrentUser())
// }

// // PublishMap publishes `contents` as a directory using the current user.
// func (a *App) PublishMap(contents map[string]Hash) (string, Hash, error) {
// 	hash, err := a.StoreMap(contents)
// 	if err != nil {
// 		return "", Hash{}, errors.Wrap(err, "PublishMap: store failed")
// 	}
//
// 	tx, err := a.PublishHash(hash)
// 	if err != nil {
// 		return "", Hash{}, errors.Wrap(err, "PublishMap: publish failed")
// 	}
//
// 	return tx, hash, nil
// }
//

// // StorePath adds `path` into ipfs, returning it's hash
// func (a *App) StorePath(path string) (Hash, error) {
// 	h, err := ipfs.Add(path)
// 	if err != nil {
// 		return Hash{}, errors.Wrap(err, "store-path: ipfs add failed")
// 	}
//
// 	return Hash{h}, nil
// }
