package dapp

import (
	"bytes"
	"flag"

	"github.com/jbenet/go-multihash"
)

// Hash represents a single hash in the dapp system
type Hash struct {
	multihash.Multihash
}

// Equals returns true if two hashes are equal
func (h Hash) Equals(other Hash) bool {
	return bytes.Equal(h.Multihash, other.Multihash)
}

// Bytes returns a copy the raw value of the hash
func (h Hash) Bytes() []byte {
	var ret bytes.Buffer
	ret.Write([]byte(h.Multihash))
	return ret.Bytes()
}

func (h Hash) String() string {
	return h.Multihash.B58String()
}

// Identity represents a single identity in the dapp system
type Identity interface {
	Equals(other Identity) bool
	PublicKey() string
	Verify(input []byte, signature []byte) error
	Sign(input []byte) ([]byte, error)
}

// IdentityProvider provides ids, either through parsing serialized values or
// generating random identities.
type IdentityProvider interface {
	ParseIdentity(str string) (Identity, error)
	RandomIdentity() (Identity, error)
	AnnounceIdentity(Identity) (TX, error)
	IsIdentityAnnounced(id Identity) (bool, error)
}

// KV reprents a ssytem that can perform a kv set/get in a decentralized
// manner.
type KV interface {
	Set(identity Identity, key string, value []byte) (TX, error)
	Get(identity Identity, key string) ([]byte, error)
}

// Store represents a module that can store and load filesystems
// addressed by their content.
type Store interface {
	// StorePath writes the contents of `path` to the store, return the hash that
	// addresses the contents.
	StorePath(path string) (Hash, error)

	// LoadPath writes the content (either directory or file) addressed by
	// `content` to `path`
	LoadPath(path string, content Hash) error
}

// TX represents the id of a transaction
type TX string

// CurrentUser returns the current process' identity within `app`
func CurrentUser(app string) Identity {
	return loginSessions[app]
}

// Login logs the current process into `app` as `user`, replacing any current
// session.
func Login(app string, user Identity) {
	loginSessions[app] = user
}

// Logout logs the current process out of `app`
func Logout(app string) {
	delete(loginSessions, app)
}

var dev = flag.Bool(
	"dapp.dev",
	false,
	"enables developer mode",
)

var printVersion = flag.Bool(
	"dapp.version",
	false,
	"print the current version and exit",
)

var printID = flag.Bool(
	"dapp.id",
	false,
	"print the app's id and exit",
)

var version = "devel"
var loginSessions map[string]Identity

func init() {
	loginSessions = map[string]Identity{}
}
