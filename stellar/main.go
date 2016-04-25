package stellar

import (
	"fmt"
	// "log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/stellar/go-stellar-base/horizon"
	"github.com/stellar/go-stellar-base/keypair"
)

// DefaultClient is the default horizon config
var DefaultClient = &Client{horizon.DefaultTestNetClient}

// Client connects to the stellar network
type Client struct {
	*horizon.Client
}

// Identity implements dapp.Identity
type Identity struct {
	keypair.KP
}

// AccountExists returns true if a stellar account at `aid` exists and is
// funded.
func AccountExists(horizon string, aid string) (bool, error) {
	url := fmt.Sprintf("%s/accounts/%s", horizon, aid)

	resp, err := http.Get(url)
	if err != nil {
		return false, errors.Wrap(err, "load account data failed")
	}

	return (resp.StatusCode >= 200 && resp.StatusCode < 300), nil
}

// FundAccount funds `aid` on the stellar network using the the friendbot at
// `horizon`.
func FundAccount(horizon string, aid string) error {
	exists, err := AccountExists(horizon, aid)
	if err != nil {
		return errors.Wrap(err, "identity existence check errored")
	}

	if exists {
		// TODO: use an actual error struct, embed the network passphrase of the
		// horizon server consulted.
		return errors.New("identity already funded")
	}

	url := fmt.Sprintf("%s/friendbot?addr=%s", horizon, aid)

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "friendbot error")
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		// TODO: use a better error by interpetting the horizon problem response
		return errors.New("friendbot failed")
	}

	return nil
}

// LoadAccountData returns a map of data values on `aid` from `horizon`
func LoadAccountData(
	horizon string,
	aid string,
) (ret map[string][]byte, err error) {

	url := fmt.Sprintf("%s/accounts/%s", horizon, aid)
	resp, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, "load account data failed")
		return
	}

	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		// TODO: better error
		err = errors.New("horizon: account load failed")
		return
	}

	return
}
