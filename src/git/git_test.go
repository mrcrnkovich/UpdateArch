package git

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) string {
	originRepo, _ := filepath.Abs("static/bare_repo")
	tmpDir := t.TempDir()

	t.Log("Origin Remote Repo location: " + originRepo)

	initCmd := exec.Command("git", "init")
	initCmd.Dir = tmpDir

	if e := initCmd.Run(); e != nil {
		t.Fatal("Couldn't build testing directory")
	}

	addRemoteOriginCmd := exec.Command(
		"/usr/bin/git",
		"remote",
		"add",
		"origin",
		originRepo,
	)
	addRemoteOriginCmd.Dir = tmpDir

	if e := addRemoteOriginCmd.Run(); e != nil {
		t.Fatal("couldn't set remote origin for git testing directory")
	}

	return tmpDir
}

func TestPull(t *testing.T) {

	/*Setup Env
	* Need a remote "origin" git repo
	*
	 */
	getTestDir := setup(t)

	// Run Pull
	testpkg := GitDir{getTestDir, false}
	if e := testpkg.Pull("origin", "master"); e != nil {
		t.Errorf("Unexpected Error while running pull, %q", e)
	}

	// Assert GitDir.Updated == true
	if testpkg.Updated != true {
		t.Error("GitDir.Updated != True ")
	}

	// Test Files Downloaded
	gotFiles, _ := filepath.Glob(filepath.Join(getTestDir, "*.txt"))
	wantFiles := []string{"test.txt"}

	for idx, gotFile := range gotFiles {
		gotFilename := filepath.Base(gotFile)
		if gotFilename != wantFiles[idx] {
			t.Errorf(" Missing Git file want: %s got: %s", wantFiles[idx], gotFilename)
		}
	}

	// Run Pull on an already up to date GitDir.
	if e := testpkg.Pull("origin", "master"); e != nil {
		t.Errorf("Unexpected Error while running pull, %q", e)
	}

	// Assert GitDir.Updated == false
	if testpkg.Updated == true {
		t.Error("Want: GitDir.Updated == false ")
	}

	// What happens when I pass it a bad path?

	// What Happens when there is no remote upstream to sync?
}
