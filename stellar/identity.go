package stellar

// PublicKey implements dapp.Identity
func (i *Identity) PublicKey() string {
	return i.KP.Address()
}

// String implements Stringer
func (i *Identity) String() string {
	return i.KP.Address()
}
