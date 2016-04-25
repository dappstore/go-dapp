package dapp

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jbenet/go-multihash"
	"github.com/pkg/errors"
	"github.com/stellar/go-stellar-base/horizon"
)

// App represents the identity for an application that is deployed using dapp
type App struct {
	ID string

	Providers struct {
		IdentityProvider
		KV
		Store
	}

	policies []Policy
}

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

// Policy values represent a policy that can change state on the app
type Policy interface {
	ApplyDappPolicy(*App) error
}

// Store represents a module that can store and load filesystems
// addressed by their content.
type Store interface {
	StoreLocalDir(path string) (Hash, error)
	LoadLocalDir(path string, content Hash) error
	NewTempDir() (string, error)
}

// TX represents the id of a transaction
type TX string

// CurrentUser returns the current process' identity within `app`
func CurrentUser(app string) Identity {
	return loginSessions[app]
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
	if *printID {
		fmt.Println(id)
		os.Exit(0)
	}

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return
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

// ManifestHash resolves the multihash for the manifest of application `id`
// using `horizons`.
func ManifestHash(id string, horizons ...string) (multihash.Multihash, error) {
	var err error

	if len(horizons) == 0 {
		return multihash.Multihash(""),
			errors.New("no verification servers specified")
	}
	manifestHashes := make([]multihash.Multihash, len(horizons))

	for i, server := range horizons {
		manifestHashes[i], err = loadRelease(server, id)
		if err != nil {
			return nil, errors.Wrap(err, "load manifest hash failed")
		}
	}

	// TODO: ensure trust threshold is satisfied

	// return the first non-nil
	for _, h := range manifestHashes {
		if h != nil {
			return h, nil
		}
	}

	// TODO: use an error struct
	return nil, errors.New("could not load any manifest hashes")
}

// PublishHash publishes `hash` using `identity`.
// func PublishHash(hash Hash, identity Identity) (string, error) {
// 	return stellar.PublishHash(
// 		defaultHorizon,
// 		identity.(*stellar.Identity),
// 		hash.Multihash,
// 	)
// }

// SetDefaultHorizon sets the default horizon server
func SetDefaultHorizon(addy string) {
	defaultHorizon = addy
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
var defaultHorizon = horizon.DefaultTestNetClient.URL
var loginSessions map[string]Identity

type manifestDisagreementError struct{}

func init() {
	loginSessions = map[string]Identity{}
}

func identityExists(id Identity) (bool, error) {
	url := fmt.Sprintf("%s/accounts/%s", defaultHorizon, id.PublicKey())

	resp, err := http.Get(url)
	if err != nil {
		return false, errors.Wrap(err, "load account data failed")
	}

	return (resp.StatusCode >= 200 && resp.StatusCode < 300), nil
}

func loadRelease(id string, server string) (multihash.Multihash, error) {
	hash, err := loadIdentityData(server, id, "dapp-release")
	if err != nil {
		return nil, errors.Wrap(err, "read identity data failed")
	}

	return hash, nil
}

func loadIdentityData(server, id, key string) ([]byte, error) {
	url := fmt.Sprintf("%s/accounts/%s/data/%s", server, id, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	req.Header.Add("Accept", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request errored")
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		// TODO: use a better error by interpetting the horizon problem response
		return nil, errors.New("request failed")
	}

	hash, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	return hash, nil
}

func loadIdentityMultihash(
	server,
	id,
	key string,
) (ret multihash.Multihash, err error) {

	ret, err = loadIdentityData(server, id, key)
	if err != nil {
		err = errors.Wrap(err, "read identity data failed")
		return
	}

	return
}

// func verifyPublication(server, id, path string) (bool, error) {
// 	publishedHash, err := loadIdentityMultihash(server, id, "dapp-publications")
// 	if err != nil {
// 		return false, errors.Wrap(err, "get publication hash failed")
// 	}
//
// 	exists, err := ipfs.Exists(publishedHash, id)
// 	if err != nil {
// 		return false, errors.Wrap(err, "directory verification failed")
// 	}
//
// 	//TODO: load the manifest, verify signatures against binaries
//
// 	return exists, nil
// }

// // ModifyDir loads an ipfs dir, modifies it according to `fn` and
// // commits it back to ipfs, returning the hash
// func ModifyDir(
// 	base multihash.Multihash,
// 	dir string,
// 	fn DirModifier,
// ) (ret multihash.Multihash, err error) {
// 	exists, err := ipfs.Exists(base, dir)
// 	if err != nil {
// 		err = errors.Wrap(err, "failed to check ipfs existence")
// 		return
// 	}
//
// 	next, err := ioutil.TempDir("", "dapp-modify-dir")
// 	if err != nil {
// 		return
// 	}
// 	defer os.RemoveAll(dir)
//
// 	if exists {
// 		err = ipfs.Get(base, dir, next)
// 		if err != nil {
// 			err = errors.Wrap(err, "failed to populate temp dir")
// 			return
// 		}
// 	}
//
// 	err = fn(next)
// 	if err != nil {
// 		err = errors.Wrap(err, "modify callback errored")
// 		return
// 	}
//
// 	ret, err = ipfs.Add(next)
// 	if err != nil {
// 		err = errors.Wrap(err, "ipfs add dailed")
// 		return
// 	}
//
// 	return
// }
