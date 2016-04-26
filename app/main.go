package dapp

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// App represents the identity for an application that is deployed using dapp
type App struct {
	ID string

	Providers struct {
		dapp.IdentityProvider
		dapp.KV
		dapp.Store
	}

	policies []Policy
}

// Policy values represent a policy that can change state on the app
type Policy interface {
	ApplyDappPolicy(*App) error
}

// NewApp creates a new dapp application with identity `id` and applies
// `policies`.
func NewApp(id string, policies ...Policy) (app *App, err error) {
	app = &App{ID: id}

	err = app.ApplyPolicies(policies...)
	if err != nil {
		err = errors.Wrap(err, "dapp: create-app failed to apply policies")
		return
	}

	if app.Providers.IdentityProvider == nil {
		err = errors.New("dapp: no identity provider initialized while applying policies")
		return
	}

	if app.Providers.KV == nil {
		err = errors.New("dapp: no kv initialized while applying policies")
		return
	}

	if app.Providers.Store == nil {
		err = errors.New("dapp: no store initialized while applying policies")
		return
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

	return
}

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
func (a *App) CurrentUser() dapp.Identity {
	return dapp.CurrentUser(a.ID)
}

// Login logs `user` into `a`
func (a *App) Login(user dapp.Identity) {
	dapp.Login(a.ID, user)
}
