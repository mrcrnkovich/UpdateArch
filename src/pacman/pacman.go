package pacman

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var writer io.Writer = os.Stdout

type CommandError struct {
	text   error
	stderr string
}

func (e CommandError) Error() string {
	return fmt.Sprintf("makepkg  %s", e.text)
}

func MakePackage(pkgPath string) ([]string, error) {

	var e error
	var pkgs []string
	if pkgs, e = packageList(pkgPath); e != nil {
		return nil, e
	}

	/*
		if err := makePackage(pkgPath); err != nil {
			return nil, err
		}
	*/

	return pkgs, nil
}

func packageList(pkgPath string) ([]string, error) {

	var result, err bytes.Buffer
	buffer := bufio.NewScanner(&result)
	PkgDestEnv := fmt.Sprintf("PKGDEST=%s", "/home/mike/.local/share/packages")

	cmd := exec.Command(
		"/usr/bin/makepkg",
		"--packagelist",
	)
	cmd.Dir = pkgPath
	cmd.Stdin = nil
	cmd.Stdout = &result
	cmd.Stderr = &err
	cmd.Env = append(cmd.Environ(), PkgDestEnv)

	if e := cmd.Run(); e != nil {
		fmt.Printf("%d", cmd.ProcessState.ExitCode())
		return nil, &CommandError{e, err.String()}
	}

	var packages []string
	for buffer.Scan() {
		txt := buffer.Text()
		packages = append(packages, txt)
		fmt.Fprintf(writer, "%s\n", txt)
	}

	return packages, nil
}

func makePackage(pkgPath string) error {

	var result, err bytes.Buffer

	cmd := exec.Command(
		"/usr/bin/makepkg",
		"--noconfirm",
		"--syncdeps",
		"--rmdeps",
		"--clean",
		"--force",
	)
	cmd.Dir = pkgPath
	cmd.Stdout = &result
	cmd.Stderr = &err
	cmd.Env = append(cmd.Environ(), "PKGDEST=/home/mike/tmp/packages")

	if e := cmd.Run(); e != nil {
		fmt.Println(err.String())
		return e
	}
	//fmt.Fprintf(writer, "%s\n", result.String())

	return nil
}

func UpdateRepo(repoPath string, pkgName string) error {

	var result, err bytes.Buffer
	fmt.Fprintf(writer, "Adding package: %s to %s\n", pkgName, repoPath)

	cmd := exec.Command("repo-add",
		"--new",
		"--remove",
		repoPath,
		pkgName,
	)
	cmd.Stdout = &result
	cmd.Stderr = &err

	if e := cmd.Run(); e != nil {
		fmt.Fprintf(writer, "Package output prior to error: %s", result.String())
		return e
	}

	return nil
}
