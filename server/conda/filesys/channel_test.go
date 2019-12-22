package filesys

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"private-conda-repo/conda/condatypes"
)

func TestConda_CRUDPackage(t *testing.T) {
	t.Parallel()

	var assert = assert.New(t)
	repo, cleanup := newTestConda()
	defer cleanup()

	chn, err := repo.CreateChannel("test-channel")
	assert.NoError(err)

	testPkg := testPackages["perfana-0.0.6-py_0.tar.bz2"]

	file, err := os.Open(testPkg.Path)
	assert.NoError(err)
	defer func() { _ = file.Close() }()

	platform := string(condatypes.NOARCH)
	pkg, err := chn.AddPackage(file, platform, "perfana")
	assert.Error(err)

	pkg, err = chn.AddPackage(file, platform, testPkg.Filename)
	assert.NoError(err)

	meta, err := chn.GetMetaInfo()
	assert.NoError(err)

	assert.Len(meta.Packages, 1)
	assert.NotNil(meta.Packages["perfana"])

	err = chn.RemoveSinglePackage(pkg)
	assert.NoError(err)
}

func TestChannel_GetMetaInfo(t *testing.T) {
	t.Parallel()

	var assert = assert.New(t)
	repo, cleanup := newTestConda()
	defer cleanup()

	chn, err := repo.CreateChannel("get-meta-info-channel")
	assert.NoError(err)

	for _, details := range testPackages {
		f, err := os.Open(details.Path)
		assert.NoError(err)
		defer func() { _ = f.Close() }()
		_, err = chn.AddPackage(f, details.Platform, details.Filename)
		assert.NoError(err)
	}

	// both packages (copulae and perfana) are registered
	meta, err := chn.GetMetaInfo()
	assert.NoError(err)
	assert.Len(meta.Packages, 2)
	assert.EqualValues("0.4.3", *meta.Packages["copulae"].Version)
	assert.EqualValues("0.0.6", *meta.Packages["perfana"].Version)

	// Remove package updates indices correctly
	err = chn.RemoveSinglePackage(&condatypes.Package{
		Name:        "perfana",
		Version:     "0.0.6",
		BuildString: "py",
		Platform:    "noarch",
	})
	assert.NoError(err)
	meta, err = chn.GetMetaInfo()
	assert.NoError(err)
	assert.EqualValues("0.0.5", *meta.Packages["perfana"].Version)
}
