package git

import (
	"fmt"
	"os/exec"
	"regexp"
)

func main() {
	fmt.Println("vim-go")
}

type GitDir struct {
	Name    string
	Path    string
	Updated bool
}

func Pull(workDir, fileName string) string {

	fmt.Printf("Now syncing: %s\n", fileName)

	gitPackage := workDir + fileName

	// Sync git files
	output, e := exec.Command("git", "-C", gitPackage, "pull").Output()
	if e != nil {
		fmt.Printf("%s", e)
	}
	fmt.Printf("%s", output)

	notNew, _ := regexp.Match(`up to date`, []byte(output))

	//clean up repo and dirs?
	if notNew {
		return fmt.Sprintf("%s", "completed")
	} else {
		return fmt.Sprintf("%s", "new version found")
	}
}
