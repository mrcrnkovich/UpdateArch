package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

type GitDir struct {
	Path    string
	Updated bool
}

func (pkg *GitDir) Pull(arg ...string) error {

	fmt.Printf("Now syncing: %s\n", pkg.Path)

	var output, err bytes.Buffer
	cmd := exec.Command("/usr/bin/git", "pull")

	cmd.Args = append(cmd.Args, arg...)
	cmd.Dir = pkg.Path
	cmd.Stdout = &output
	cmd.Stderr = &err

	// Sync git files
	if e := cmd.Run(); e != nil {
		return fmt.Errorf("%s %s", err.String(), e)
	}
	//clean up repo and dirs?
	notNew, _ := regexp.Match(`up to date`, output.Bytes())
	pkg.Updated = !notNew

	return nil
}

func Clone(workDir, fileName string) {
	return
}
