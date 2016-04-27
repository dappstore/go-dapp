package claim

var _ MakesClaims = &MockClaimer{}
var _ Verifier = &MockVerifier{}

// MockClaimer is a mock that implements MakesClaims
type MockClaimer struct {
	Name     string
	Identity string
	Claims   string
}

// ClaimerName implements `MakesClaims`
func (t *MockClaimer) ClaimerName() string { return t.Name }

// ClaimerIdentity implements `MakesClaims`
func (t *MockClaimer) ClaimerIdentity() string { return t.Identity }

// ClaimerClaims implements `MakesClaims`
func (t *MockClaimer) ClaimerClaims() string { return t.Claims }

// MockVerifier is a mock that implements Verifier
type MockVerifier struct {
	fn func(*Claims) error
}

// VerifyClaims implements `Verifier`
func (t *MockVerifier) VerifyClaims(c *Claims) error { return t.fn(c) }
