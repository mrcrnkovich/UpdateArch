package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"git"
	"pacman"
)

var writer io.Writer = os.Stdout
var AurPath string
var RepoPath string

var ForceBuild bool

func init() {

	HOMEDIR, _ := os.UserHomeDir()

	if val, ok := os.LookupEnv("AUR_PATH"); ok {
		AurPath = val
	} else {
		AurPath = HOMEDIR + "/.local/share/aur/"
	}
	if val, ok := os.LookupEnv("REPO_PATH"); ok {
		RepoPath = val
	} else {
		RepoPath = HOMEDIR + "/.local/share/packages/mcrnkovich.db.tar.gz"
	}

	flag.BoolVar(&ForceBuild, "force", false, "force building of packages")
	flag.Parse()

	log.Println(AurPath)
	log.Println(RepoPath)
}

func main() {
	os.Setenv("PKGDEST", "/home/packages")
	defer os.Unsetenv("PKGDEST")

	files, e := os.ReadDir(AurPath)
	if e != nil {
		log.Fatal("Unable to read dir:", AurPath)
	}

	for _, file := range files {
		UpdateAndSync(file.Name())
	}
}

// if file point to a Dir
func UpdateAndSync(fileName string) {

	GitPkg := git.GitDir{fileName, AurPath, false}
	gitPack, e := git.Pull(&GitPkg)

	if e != nil {
		fmt.Fprintf(writer, "%s", e)
	} else if !gitPack.Updated && !ForceBuild {
		fmt.Fprintf(writer, "No updates found for package: %s\n", fileName)
		return
	}

	packages, err := pacman.MakePackage(AurPath + fileName)
	if err != nil {
		log.Printf("%s", err)
	}

	for _, pkg := range packages {
		if e := pacman.UpdateRepo(RepoPath, pkg); e != nil {
			fmt.Fprintf(writer, "%s", e)
		}
	}
}
