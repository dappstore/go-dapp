package claim

var _ MakesClaims = &MockClaimer{}

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
