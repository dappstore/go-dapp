package hash

import (
	"github.com/dappstore/go-dapp"
)

// Hasher reprents a module that can hash a local directory to a
type Hasher interface {
	HashLocalPath(path string) dapp.Hash
}
