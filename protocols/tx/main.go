package tx

import (
	"github.com/dappstore/go-dapp"
)

// System represents a module that can commit transactions and then check
// whether a transaction was committed some time in the past.
type System interface {
	Commit(
		source dapp.Identity,
		tx dapp.TX,
		signers []dapp.Identity,
	) (dapp.Hash, error)

	Committed(hash dapp.Hash) (bool, error)
}
