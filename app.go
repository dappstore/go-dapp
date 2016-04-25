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

// Login logs `user` into `a`
func (a *App) Login(user Identity) {
	Login(a.ID, user)
}
