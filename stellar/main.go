package stellar

import (
	"fmt"
	// "log"
	"encoding/base64"
	"encoding/json"
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
func AccountExists(h *horizon.Client, aid string) (bool, error) {
	url := fmt.Sprintf("%s/accounts/%s", h.URL, aid)

	resp, err := http.Get(url)
	if err != nil {
		return false, errors.Wrap(err, "load account data failed")
	}

	return (resp.StatusCode >= 200 && resp.StatusCode < 300), nil
}

// FundAccount funds `aid` on the stellar network using the the friendbot at
// `horizon`.
func FundAccount(h *horizon.Client, aid string) (string, error) {
	exists, err := AccountExists(h, aid)
	if err != nil {
		return "", errors.Wrap(err, "identity existence check errored")
	}

	if exists {
		// TODO: use an actual error struct, embed the network passphrase of the
		// horizon server consulted.
		return "", errors.New("identity already funded")
	}

	url := fmt.Sprintf("%s/friendbot?addr=%s", h.URL, aid)

	var result struct {
		Hash string `json:"hash"`
	}

	err = decodeGet(url, &result)
	if err != nil {
		return "", errors.Wrap(err, "fund account: friendbot request failed")
	}

	return result.Hash, nil
}

// LoadAccountData returns a map of data values on `aid` from `horizon`
func LoadAccountData(
	h *horizon.Client,
	aid string,
) (ret map[string][]byte, err error) {

	url := fmt.Sprintf("%s/accounts/%s", h.URL, aid)

	var result struct {
		Data map[string]string `json:"data"`
	}

	err = decodeGet(url, &result)
	if err != nil {
		err = errors.Wrap(err, "load account data: hoirzon request failed")
		return
	}

	ret = map[string][]byte{}
	for k, v := range result.Data {
		ret[k], err = base64.StdEncoding.DecodeString(v)
		if err != nil {
			err = errors.Wrap(err, "load account data: hoirzon request failed")
			return
		}
	}

	return
}

func decodeGet(url string, dest interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "horizon: request errored")
	}

	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return errors.New("horizon: request failed")
	}

	enc := json.NewDecoder(resp.Body)
	err = enc.Decode(dest)
	if err != nil {
		return errors.Wrap(err, "horizon: decode response failed")
	}

	return nil
}
