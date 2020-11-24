package resolve_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x-plugins/jx-mink/pkg/cmd/resolve"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner/fakerunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/stretchr/testify/require"
)

func TestMinkResolve(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "could not create temp dir")

	t.Logf("running tests in %s\n", tmpDir)

	fs, err := ioutil.ReadDir("test_data")

	for _, f := range fs {
		if !f.IsDir() {
			continue
		}
		name := f.Name()
		srcFile := filepath.Join("test_data", name)
		require.DirExists(t, srcFile)

		destDir := filepath.Join(tmpDir, name)
		err = files.CopyDirOverwrite(srcFile, destDir)
		require.NoError(t, err, "failed to copy %s to %s", srcFile, destDir)

		runner := &fakerunner.FakeRunner{}

		_, o := resolve.NewCmdMinkResolve()
		o.Dir = destDir
		o.CommandRunner = runner.Run
		err = o.Run()
		require.NoError(t, err, "failed for test %s", name)

		for _, c := range runner.OrderedCommands {
			t.Logf("test %s invoked: %s\n", name, c.CLI())
		}

		if name != "no-image" {
			require.NotEmpty(t, runner.OrderedCommands, "should have ran a command for %s", name)
		}

	}
}

func TestMinkResolveLocalKaniko(t *testing.T) {
	name := "local-kaniko"

	runner := &fakerunner.FakeRunner{}

	_, o := resolve.NewCmdMinkResolve()
	o.Dir = filepath.Join("test_data", "dockerfile")
	o.CommandRunner = runner.Run
	o.PlainKaniko = true
	err := o.Run()
	require.NoError(t, err, "failed for test %s", name)

	for _, c := range runner.OrderedCommands {
		t.Logf("test %s invoked: %s\n", name, c.CLI())
	}
	require.NotEmpty(t, runner.OrderedCommands, "should have ran a command for %s", name)
}
