package module_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/nobe4/gtfo/internal/module"
	"github.com/stretchr/testify/assert"
)

const (
	moduleExpected = "bit.ly/3dKQQSt"
)

// prepare will create a fake directory to create (or not) a go.mod in.
func prepare(t *testing.T, withMod bool) func() {
	dir, err := ioutil.TempDir("", "example")
	assert.NoError(t, err)

	if withMod {
		content := fmt.Sprintf("module %s\n\ngo 3.14", moduleExpected)
		tmpfn := filepath.Join(dir, "go.mod")
		err := ioutil.WriteFile(tmpfn, []byte(content), 0600)

		assert.NoError(t, err)
	}

	err = os.Chdir(dir)
	assert.NoError(t, err)

	return func() { os.RemoveAll(dir) }
}

func TestNoMod(t *testing.T) {
	cleanup := prepare(t, false)
	defer cleanup()

	moduleFound, err := module.Get()
	assert.Error(t, err)
	assert.Empty(t, moduleFound)
}

func TestWithMod(t *testing.T) {
	cleanup := prepare(t, true)
	defer cleanup()

	moduleFound, err := module.Get()
	assert.NoError(t, err)
	assert.Equal(t, moduleExpected, moduleFound)
}
