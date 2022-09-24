package pacman

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("vim-go")
}

func MakePackage(gitPackage string) (string, error) {

	e := os.Chdir(gitPackage)
	if e != nil {
		fmt.Printf("Error changing into %s\n", e.Error())
	}

	Package, _ := exec.Command("makepkg", "--packagelist").Output()
	pkgName := strings.TrimRight(string(Package), "\n")

	fmt.Fprintf(os.Stdin, "Make package: %s\n", gitPackage)

	/*
		out, e := exec.Command(
			"makepkg",
			"--noconfirm",
			"--syncdeps",
			"--rmdeps",
			"--clean",
			"--force",
		).Output()

		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(string(out))
	*/

	return pkgName, nil
}

func UpdateRepo(repoPath string, pkgName string) error {

	fmt.Fprintf(os.Stdin, "Adding package: %s to %s\n", pkgName, repoPath)

	output, e := exec.Command("repo-add",
		"--new",
		"--remove",
		repoPath,
		pkgName,
	).Output()
	if e != nil {
		return e
	}

	fmt.Fprintf(os.Stdin, "This package was made: %s", output)

	return nil
}

func printPackageInfo(pkg []byte) {
	p := regexp.MustCompile(`(?:\/home\/packages\/)(?P<package>.+)-(?P<version>[0-9]+\.?[0-9]*\.?[0-9]*)-(?P<q>\d)-(?P<arch>.+)(?:\.pkg\.tar\.zst)`)
	matches := p.FindSubmatch(pkg)
	keys := p.SubexpNames()

	for i, m := range matches {
		if i > 0 {
			fmt.Printf("%-15s: %q\n", keys[i], m)
		}
	}
}
