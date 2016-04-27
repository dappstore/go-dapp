package app

var _ Policy = &IdentityProvider{}
var _ Policy = &KV{}
var _ Policy = &RunVerification{}
var _ Policy = &Store{}
var _ Policy = &VerifySelf{}
var _ Policy = Name("me")
var _ Policy = Developer("GSDSED")
var _ Policy = Description("It just spins")
