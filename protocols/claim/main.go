// Package claim implements the claim protocol.
//
// As part of applying policies within a process, claims can be made.  these
// claims can be serialized to json and recorded. claims can be trusted by
// identities.
//
// A claim protocol is one of the lowest level protocols in dapp.  A claim
// represents an unprovable statement that can be interpretted by both computers
// and humans.  Trust in a claim is established through any trust protocol.
package claim

import (
	"os"
	"sync"

	"github.com/Jeffail/gabs"
	"github.com/spf13/afero"
)

// ClaimersClaimPath is the claim path for the set of currently registers claimers.
const ClaimersClaimPath = "protocols.claim.claimers"

// LockerClaimersClaimPath is the claim path for the set of claimers registered
// at the time the claim.LockClaimers() protocol was run.
const LockerClaimersClaimPath = "protocols.claim.locked-claimers"

// Default represents the default instantiation of this protocol
var Default *Protocol

// Claims represents a tree of claims.
type Claims struct {
	data *gabs.Container
}

// MakesClaims represents a type that makes claims as part of this protocol.  A
// claimer can be queried for it's local name (a friendly name for humans), an
// identity, and a set of claims that it makes about itself.
type MakesClaims interface {
	ClaimerName() string
	ClaimerIdentity() string
	ClaimerClaims() string
}

// Protocol represents a single set of claims made by the running application.
type Protocol struct {
	Claims *Claims

	lock           sync.Mutex
	claimersLocked bool
	claimers       []MakesClaims
}

// New creates a new claim protocol
func New() *Protocol {
	return &Protocol{
		Claims: NewClaims(),
	}
}

// NewClaims creates a new claims struct
func NewClaims() *Claims {
	return &Claims{data: gabs.New()}
}

// AddClaimer registers `c` as an identity that can make claims
// on the default instance of this protocol`.
func AddClaimer(c MakesClaims) error {
	return Default.AddClaimer(c)
}

// CurrentClaims returns the claims that have been recorded on the default
// instance of this protocol.
func CurrentClaims() string {
	return Default.CurrentClaims()
}

// LockClaimers prevents further claimers from being added to the default
// instance of this`gq protocol
func LockClaimers() error {
	return Default.LockClaimers()
}

// Make makes a claim on the default claim protocol
func Make(path string, value interface{}) error {
	return Default.Make(path, value)
}

// Push pushes a claim on the default claim protocol
func Push(path string, value interface{}) error {
	return Default.Make(path, value)
}

// WriteFile writes the claims made on the default claim protocol to disk.
func WriteFile(fs afero.Fs, path string, perm os.FileMode) error {
	return Default.WriteFile(fs, path, perm)
}

func init() {
	Default = New()
}
