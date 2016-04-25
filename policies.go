package dapp

import (
	"github.com/dappstore/agree"
)

// AgreementOracle represents the dapp policy that when applied to the current
// process causes dapp consider using `Horizon` when calculating agreement.
type AgreementOracle struct {
	agree.Oracle
}

var _ Policy = &AgreementOracle{}

// ApplyDappPolicy applies `p` to `app`
func (p *AgreementOracle) ApplyDappPolicy(app *App) error {
	// app.agreements.Oracles = append(app.agreements.Oracles, p.Oracle)
	return nil
}

// AgreementPolicy represents the dapp policy that adds policy to the current
// process' agreement system.
type AgreementPolicy struct {
	agree.Policy
}

var _ Policy = &AgreementPolicy{}

// ApplyDappPolicy applies `p` to `app`
func (p *AgreementPolicy) ApplyDappPolicy(app *App) error {
	// app.agreements.Policies = append(app.agreements.Policies, p.Policy)
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
