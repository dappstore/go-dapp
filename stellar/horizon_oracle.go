package stellar

import (
	"github.com/dappstore/agree"
	"github.com/pkg/errors"
)

var _ agree.Oracle = HorizonOracle("https://horizon.stellar.org")

// GetOracleView implements `agree.Oracle`
func (o HorizonOracle) GetOracleView(domain, key string) ([]byte, error) {
	data, err := LoadAccountData(string(o), domain)
	if err != nil {
		return nil, errors.Wrap(err, "load account data failed")
	}

	return data[key], nil
}
