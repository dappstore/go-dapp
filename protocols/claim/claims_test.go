package claim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaims_String(t *testing.T) {
	var claims *Claims
	// nil claims
	assert.Equal(t, "{}", claims.String())

	// data is nil panics
	claims = &Claims{}
	assert.Panics(t, func() { _ = claims.String() })

	// populated claims
	claims = NewClaims()
	claims.make("foo", 1)
	assert.Equal(t, `{"foo":1}`, claims.String())
}

func TestClaims_make(t *testing.T) {
	c := NewClaims()
	data := c.data

	// root claim
	if assert.NoError(t, c.make("foo", 1)) {
		assert.Equal(t, 1, data.Path("foo").Data())
	}

	// nested claim
	if assert.NoError(t, c.make("nested.1.2", 1)) {
		assert.Equal(t, 1, data.Path("nested.1.2").Data())
	}

	// reclaim fails
	assert.Error(t, c.make("foo", "something else"))
	assert.Equal(t, 1, data.Path("foo").Data())

	// fails on set
	require.NoError(t, c.push("blah", 1))
	assert.Error(t, c.make("blah", "something else"))
}

func TestClaims_push(t *testing.T) {
	c := NewClaims()
	data := c.data

	// root claim
	if assert.NoError(t, c.push("foo", 1)) {
		assert.Equal(t, "[1]", data.Path("foo").String())
	}

	// nested claim
	if assert.NoError(t, c.push("nested.1.2", 1)) {
		assert.Equal(t, "[1]", data.Path("nested.1.2").String())
	}

	// fails on non-set claim
	require.NoError(t, c.make("blah", 1))
	assert.Error(t, c.push("blah", "something else"))
	t.Log(data.String())
}
