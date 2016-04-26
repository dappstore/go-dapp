package publish

import (
	"github.com/dappstore/go-dapp"
	"github.com/dappstore/go-dapp/protocols/claim"
	"github.com/dappstore/go-dapp/protocols/dfs"
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
) (tx dapp.TX, publication dapp.Hash, err error) {

	pdfs := dfs.New(sys.store)
	claims, err := pdfs.StoreString(claim.Default.Claims.String())

	// merge current processe's claims file into hash
	publication, err = pdfs.MergeAtPath(contents, ".dapp/claims/publish", claims)
	if err != nil {
		err = errors.Wrap(err, "protocol-publish: failed to merge claims with content")
		return
	}

	tx, err = sys.kv.Set(publisher, "dapp:publications", contents.Bytes())
	if err != nil {
		err = errors.Wrap(err, "protocol-publish: failed to set publication hash")
		return
	}

	return
}
