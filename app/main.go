package app

import (
	"fmt"
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
	"sync"
)

// App represents the identity for an application that is deployed using dapp
type App struct {
	ID string

	Providers struct {
		dapp.IdentityProvider
		dapp.KV
		dapp.Store
	}

	once     sync.Once
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
	app.once.Do(func() {
		err = app.init(policies)
	})

	return
}

// NewPolicy creates a new composite policy.
func NewPolicy(name string, policies ...Policy) Policy {
	return &compositePolicy{name, policies}
}

type compositePolicy struct {
	name     string
	policies []Policy
}

// ApplyDappPolicy implements `Policy`
func (p *compositePolicy) ApplyDappPolicy(app *App) error {
	for i, cp := range p.policies {
		err := cp.ApplyDappPolicy(app)
		if err != nil {
			msg := fmt.Sprintf("%s: child-policy %d failed", p.name, i)
			return errors.Wrap(err, msg)
		}
	}

	return nil
}

// fnPolicy is a helper to make it easy to build policies from functions
type fnPolicy struct {
	name string
	fn   func(app *App) error
}

// ApplyDappPolicy implements `Policy`
func (p *fnPolicy) ApplyDappPolicy(app *App) error {
	err := p.fn(app)
	if err != nil {
		msg := fmt.Sprintf("%s: policy failed", p.name)
		return errors.Wrap(err, msg)
	}

	return nil
}
