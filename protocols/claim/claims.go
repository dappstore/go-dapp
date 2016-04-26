package claim

import (
	"github.com/pkg/errors"
)

func (c *Claims) String() string {
	if c == nil {
		return "{}"
	}

	if c.data == nil {
		panic("Invalid claims struct. data is nil")
	}

	return c.data.String()
}

func (c *Claims) make(path string, value interface{}) error {
	data := c.data

	// if a value exists at the path, error
	if data.ExistsP(path) {
		return errors.New("claim already set")
	}

	data.SetP(value, path)
	return nil
}

func (c *Claims) push(path string, value interface{}) error {
	var err error
	data := c.data

	if !data.ExistsP(path) {
		_, err = data.ArrayP(path)
		if err != nil {
			return errors.Wrap(err, "claim at path is not an array")
		}
	}

	err = data.ArrayAppendP(value, path)
	if err != nil {
		return errors.Wrap(err, "claim at path is not an array")
	}

	return nil
}
