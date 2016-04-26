package claim

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMake(t *testing.T) {
	p := New()

	// root claim
	if assert.NoError(t, p.Make("foo", 1)) {
		assert.Equal(t, 1, p.Claims.Path("foo").Data())
	}

	// nested claim
	if assert.NoError(t, p.Make("nested/1/2", 1)) {
		assert.Equal(t, 1, p.Claims.Path("nested/1/2").Data())
	}

	// reclaim fails
	assert.Error(t, p.Make("foo", "something else"))
	assert.Equal(t, 1, p.Claims.Path("foo").Data())

	// fails on set
	require.NoError(t, p.Push("blah", 1))
	assert.Error(t, p.Make("blah", "something else"))
}

func TestPush(t *testing.T) {
	p := New()

	// root claim
	if assert.NoError(t, p.Push("foo", 1)) {
		assert.Equal(t, "[1]", p.Claims.Path("foo").String())
	}

	// nested claim
	if assert.NoError(t, p.Push("nested/1/2", 1)) {
		assert.Equal(t, "[1]", p.Claims.Path("nested/1/2").String())
	}

	// fails on non-set claim
	require.NoError(t, p.Make("blah", 1))
	assert.Error(t, p.Push("blah", "something else"))
	t.Log(p.Claims.String())
}

func TestWriteFile(t *testing.T) {
	p := New()
	require.NoError(t, p.Make("foo", 1))

	var fs afero.Fs
	fs = afero.NewMemMapFs()

	// data is correct
	assert.NoError(t, p.WriteFile(fs, "claim", 0700))
	data, err := afero.ReadFile(fs, "claim")
	if assert.NoError(t, err) {
		assert.Equal(t, `{"foo":1}`, string(data))
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
