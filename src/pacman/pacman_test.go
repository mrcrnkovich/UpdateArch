package pacman

import (
	"bytes"
	"os"
	"testing"
)

func TestMakeRepo(t *testing.T) {

	tmpdir := t.TempDir()
	t.Setenv("PKGDEST", tmpdir)
	cwd, _ := os.Getwd()
	t.Log(cwd)

	writer = new(bytes.Buffer)

	var tests = []struct {
		input string
		want  []string
	}{
		{
			cwd + "/static/single",
			[]string{
				tmpdir + "/single-test-1.4.0-1-x86_64.pkg.tar.zst",
			},
		},
		{
			cwd + "/static/multi",
			[]string{
				tmpdir + "/multi-first-1.4.0-1-x86_64.pkg.tar.zst",
				tmpdir + "/multi-second-1.4.0-1-x86_64.pkg.tar.zst",
				tmpdir + "/multi-third-1.4.0-1-x86_64.pkg.tar.zst",
			},
		},
	}

	for _, pkg := range tests {

		t.Log(pkg)
		pkgName, e := MakePackage(pkg.input)
		if e != nil {
			t.Error(e)
		}

		if len(pkg.want) != len(pkgName) {
			t.Error("Wanted Pkg len, got Pkg len", len(pkg.want), len(pkgName))
		} else {
			for i, result := range pkgName {

				// return result package name matches expected result
				if result != pkg.want[i] {
					t.Error("Wanted, got", pkg.want[i], result)
				} else {
					t.Logf("%s: Successful match: %s\n", t.Name(), result)
				}

				// pkg was created built
				_, e := os.Stat(result)
				if os.IsNotExist(e) {
					t.Error(" Failed to build makepkg ")
				}
			}
		}
	}

	// Check package was actually build
}

func TestUpdateRepo(t *testing.T) {
	tmpDir := t.TempDir()

	repoPath := tmpDir + "test.db.tar.gz"
	pkgName := "/home/packages/wlsunset-0.2.0-2-x86_64.pkg.tar.zst"
	if e := UpdateRepo(repoPath, pkgName); e != nil {
		t.Error("found error", e)
	} else {
		t.Logf("%s: successfully built Repo\n", t.Name())
	}
}
