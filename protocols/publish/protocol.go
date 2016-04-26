package publish

import (
	"github.com/dappstore/go-dapp"
	"github.com/pkg/errors"
)

// GetPublications resolves the publisher's latest published code
func (sys *Protocol) GetPublications(
	publisher dapp.Identity,
) (hash dapp.Hash, err error) {
	kv := sys.kv
	bytes, err := kv.Get(publisher, "dapp:publications")
	if err != nil {
		err = errors.Wrap(err, "protocol-publish: failed to get publication hash")
		return
	}

	hash.Multihash = bytes
	return
}

// SetPublications overwrites the publisher's publications hash using the hash
// for the local directory at `path`.
func (sys *Protocol) SetPublications(
	publisher dapp.Identity,
	contents dapp.Hash,
) (tx dapp.TX, hash dapp.Hash, err error) {
	kv := sys.kv

	tx, err = kv.Set(publisher, "dapp:publications", contents.Bytes())
	if err != nil {
		err = errors.Wrap(err, "protocol-publish: failed to set publication hash")
		return
	}

	hash = contents

	return
}