package publish

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// SetPublications overwrites the publisher's publications hash using the hash
// for the local directory at `path`.
func (sys *System) SetPublications(
	publisher dapp.Identity,
	contents dapp.Hash,
) (tx dapp.TX, hash dapp.Hash, err error) {

	tx, err = sys.App.KV().Set(publisher, "dapp:publications", contents.Bytes())
	if err != nil {
		err = errors.Wrap(err, "protocol-publish: failed to set publication hash")
		return
	}

	hash = contents

	return
}
