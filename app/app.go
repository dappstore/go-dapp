package app

import (
	"github.com/dappstore/go-dapp"
	"github.com/dappstore/go-dapp/protocols/claim"
	"github.com/pkg/errors"
)

// ApplyPolicy applies `p` t `a`
func (a *App) ApplyPolicy(p Policy) error {
	if p == nil {
		return errors.New("policy is nil")
	}

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
func (a *App) CurrentUser() dapp.Identity {
	return dapp.CurrentUser(a.ID)
}

// Login logs `user` into `a`
func (a *App) Login(user dapp.Identity) {
	dapp.Login(a.ID, user)
}

// SendPayment sends a simple payment.
//
// NOTE: this is not intended to be the final api... it's just a prototype
func (a *App) SendPayment(dest, amount string) (dapp.TX, error) {
	// watch for payments to developer address.  When one is seen, return

	return dapp.TX(""), errors.New("not implemented")
}

// WaitForPayment waits for a developer to publish
//
// NOTE: this is not intended to be the final api... it's just a prototype
func (a *App) WaitForPayment() (dapp.TX, error) {
	// watch for payments to developer address.  When one is seen, return

	return dapp.TX(""), errors.New("not implemented")
}

func (a *App) init(policies []Policy) error {
	err := a.ApplyPolicies(policies...)
	if err != nil {
		return errors.Wrap(err, "dapp: create-app failed to apply policies")
	}

	if a.Providers.IdentityProvider == nil {
		return errors.New("dapp: no identity provider initialized while applying policies")
	}

	if a.Providers.KV == nil {
		return errors.New("dapp: no kv initialized while applying policies")
	}

	if a.Providers.Store == nil {
		return errors.New("dapp: no store initialized while applying policies")
	}

	err = claim.LockClaimers()
	if err != nil {
		return errors.Wrap(err, "dapp: failed to lock claimers")
	}

	// TODO: extract these to policies
	// if *printID {
	// 	fmt.Println(id)
	// 	os.Exit(0)
	// }
	//
	// if *printVersion {
	// 	fmt.Println(version)
	// 	os.Exit(0)
	// }

	return nil
}
