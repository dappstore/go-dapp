package claim

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddClaimer(t *testing.T) {
	p := New()
	data := p.Claims.data

	// add a claimer
	err := p.AddClaimer(&MockClaimer{
		Name:     "Mocky",
		Identity: "GCG26FSCQEVSHQHUUHPMZQKIB76CIURZSPZ2QXEPGPQSN6JMO3WXXIQM",
	})
	assert.NoError(t, err)
	assert.Equal(t,
		`[{"Name":"Mocky","Identity":"GCG26FSCQEVSHQHUUHPMZQKIB76CIURZSPZ2QXEPGPQSN6JMO3WXXIQM","Claims":""}]`,
		data.Path(ClaimersClaimPath).String())

	// add a second claimer
	err = p.AddClaimer(&MockClaimer{
		Name:     "Mocky2: Electric Boogaloo",
		Identity: "GB3YVHSOJINX357I6FKK4K22SXIPTNGAW7GZIOI54DPLUICKNARMPAAW",
	})
	assert.NoError(t, err)
	assert.Len(t, data.Path(ClaimersClaimPath).Data(), 2)

	// lock claimer
	err = p.LockClaimers()
	require.NoError(t, err)

	// add fails
	err = p.AddClaimer(&MockClaimer{
		Name:     "Mocky3: Now with more mocking",
		Identity: "GBPSRJEFPCHHSMRVUKXLJTQ2PHUOQ5SOTQN2XPUJQFIZSMUPRWRYHOPJ",
	})
	assert.Error(t, err)

	// TODO: prevent name conflicts and prevent duplicate identity
	// TODO: ensure claimer is trusted
}

func TestLockClaimers(t *testing.T) {

	// locked with no claimers
	p := New()
	assert.NoError(t, p.LockClaimers())
	assert.Equal(t,
		`{}`,
		p.Claims.data.Path(LockerClaimersClaimPath).Data().(string),
	)
	// locked after one claimer
	p = New()
	require.NoError(t, p.AddClaimer(&MockClaimer{
		Name:     "Mocky",
		Identity: "GCG26FSCQEVSHQHUUHPMZQKIB76CIURZSPZ2QXEPGPQSN6JMO3WXXIQM",
	}))
	assert.NoError(t, p.LockClaimers())
	assert.Equal(t,
		`[{"Name":"Mocky","Identity":"GCG26FSCQEVSHQHUUHPMZQKIB76CIURZSPZ2QXEPGPQSN6JMO3WXXIQM","Claims":""}]`,
		p.Claims.data.Path(LockerClaimersClaimPath).Data().(string),
	)

	// relocking fails
	p = New()
	assert.NoError(t, p.LockClaimers())
	assert.Error(t, p.LockClaimers())
}

func TestMake(t *testing.T) {
	p := New()
	data := p.Claims.data

	// NOTE: even though the behavior below is tested as part of the Claims
	// struct, we re-test it here since the Claims structs methods are all package
	// internal.

	// root claim
	if assert.NoError(t, p.Make("foo", 1)) {
		assert.Equal(t, 1, data.Path("foo").Data())
	}

	// nested claim
	if assert.NoError(t, p.Make("nested.1.2", 1)) {
		assert.Equal(t, 1, data.Path("nested.1.2").Data())
	}

	// reclaim fails
	assert.Error(t, p.Make("foo", "something else"))
	assert.Equal(t, 1, data.Path("foo").Data())

	// fails on set
	require.NoError(t, p.Push("blah", 1))
	assert.Error(t, p.Make("blah", "something else"))
}

func TestPush(t *testing.T) {
	p := New()
	data := p.Claims.data

	// NOTE: even though the behavior below is tested as part of the Claims
	// struct, we re-test it here since the Claims structs methods are all package
	// internal.

	// root claim
	if assert.NoError(t, p.Push("foo", 1)) {
		assert.Equal(t, "[1]", data.Path("foo").String())
	}

	// nested claim
	if assert.NoError(t, p.Push("nested.1.2", 1)) {
		assert.Equal(t, "[1]", data.Path("nested.1.2").String())
	}

	// fails on non-set claim
	require.NoError(t, p.Make("blah", 1))
	assert.Error(t, p.Push("blah", "something else"))
	t.Log(data.String())
}

func TestWriteFile(t *testing.T) {
	p := New()

	require.NoError(t, p.Make("foo", 1))

	var fs afero.Fs
	fs = afero.NewMemMapFs()

	// data is correct
	assert.NoError(t, p.WriteFile(fs, "claim", 0700))
	read, err := afero.ReadFile(fs, "claim")
	if assert.NoError(t, err) {
		assert.Equal(t, `{"foo":1}`, string(read))
	}

	osfs := afero.NewOsFs()
	dir := afero.GetTempDir(osfs, "claim-tests")
	t.Log("dir:", dir)
	fs = afero.NewBasePathFs(osfs, dir)

	// permissions are recorded properly
	assert.NoError(t, p.WriteFile(fs, "claim", 0700))
	stat, err := fs.Stat("claim")
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0700), stat.Mode())

}
