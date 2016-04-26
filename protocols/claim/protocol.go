// Package claim implements the claim protocol.
//
// As part of applying policies within a process, claims can be made.  these
// claims can be serialized to json and recorded. claims can be trusted.
package claim

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// AddClaimer registers `c` in the protocol as an identity that can make claims
// on the default instance of this protocol`.
func (p *Protocol) AddClaimer(c MakesClaims) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.claimersLocked {
		return errors.New("protocol-claim: cannot add claimer after lock")
	}

	p.claimers = append(p.claimers, c)

	var claimerClaim struct {
		Name     string
		Identity string
		Claims   string
	}

	claimerClaim.Name = c.ClaimerName()
	claimerClaim.Identity = c.ClaimerIdentity()
	claimerClaim.Claims = c.ClaimerClaims()

	err := p.Claims.push(ClaimersClaimPath, claimerClaim)
	if err != nil {
		return errors.Wrap(err, "protocol-claim: failed to push claim while adding a claimer")
	}

	return nil
}

// CurrentClaims returns the claims that have been recorded on `p`
func (p *Protocol) CurrentClaims() string {
	p.lock.Lock()
	defer p.lock.Unlock()
	data := p.Claims.data

	return data.String()
}

// LockClaimers prevents further claimers from being added to this protocol
func (p *Protocol) LockClaimers() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.claimersLocked {
		return errors.New("protocol-claim: claimers already locked")
	}

	err := p.Claims.make(
		LockerClaimersClaimPath,
		p.Claims.data.Path(ClaimersClaimPath).String(),
	)
	if err != nil {
		return errors.Wrap(err, "protocol-claim: failed to make locked-claimers claim")
	}
	p.claimersLocked = true

	return nil
}

// Make records a claim made by the current process.
func (p *Protocol) Make(path string, value interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.Claims.make(path, value)
}

// Push pushes a claim on to an array at `path` for current process.
func (p *Protocol) Push(path string, value interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.Claims.push(path, value)
}

// WriteFile saves the current process' claims to disk
func (p *Protocol) WriteFile(fs afero.Fs, path string, perm os.FileMode) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	data := p.Claims.data

	err := afero.WriteFile(fs, path, data.Bytes(), perm)
	if err != nil {
		return errors.Wrap(err, "claim write failed")
	}

	return nil
}
