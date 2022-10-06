package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

type GitDir struct {
	Name    string
	Path    string
	Updated bool
}

func Pull(pkg *GitDir) (*GitDir, error) {

	fmt.Printf("Now syncing: %s\n", pkg.Name)

	var output, err bytes.Buffer
	cmd := exec.Command("/usr/bin/git", "pull")
	cmd.Dir = pkg.Path + pkg.Name
	cmd.Stdin = nil
	cmd.Stdout = &output
	cmd.Stderr = &err

	// Sync git files
	if e := cmd.Run(); e != nil {
		return nil, e
	}
	//clean up repo and dirs?
	notNew, _ := regexp.Match(`up to date`, output.Bytes())
	pkg.Updated = !notNew

	return pkg, nil
}

func Clone(workDir, fileName string) {
	return
}
