package publish

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// SetPublications overwrites the publisher's publications hash using the hash
// for the local directory at `path`.
func (sys *System) SetPublications(
	publisher dapp.Identity,
	path string,
) (tx dapp.TX, hash dapp.Hash, err error) {

	hash, err = sys.App.Store().StoreLocalDir(path)
	if err != nil {
		err = errors.Wrap(err, "protocol-publish: failed to store local path")
		return
	}

	tx, err = sys.App.KV().Set(publisher, "dapp:publications", hash.Multihash)
	if err != nil {
		err = errors.Wrap(err, "protocol-publish: failed to set publication hash")
		return
	}

	return
}
