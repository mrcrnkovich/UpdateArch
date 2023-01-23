package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"

	"git"
	"pacman"
)

var writer io.Writer = os.Stdout
var AurPath string
var RepoPath string

var Filename string
var ForceBuild bool

func init() {

	HOMEDIR, _ := os.UserHomeDir()

	AurPath = HOMEDIR + "/.local/share/aur/"
	if val, ok := os.LookupEnv("AUR_PATH"); ok {
		AurPath = val
	}

	RepoPath = HOMEDIR + "/.local/share/packages/mcrnkovich.db.tar.gz"
	if val, ok := os.LookupEnv("REPO_PATH"); ok {
		RepoPath = val
	}

	flag.BoolVar(&ForceBuild, "force", false, "force building of packages")
	flag.StringVar(&Filename, "filename", "", "Filename for single file")

	log.Println(AurPath)
	log.Println(RepoPath)

	DependencyCheck()
}

func DependencyCheck() {
	// Dependency Check
	// Confirm Git installed and Pacman installed
	if _, e := exec.LookPath("git"); e != nil {
		log.Fatal("Must have git installed and included in the PATH")
	}
	if _, e := exec.LookPath("pacman"); e != nil {
		log.Fatal("Must have pacman installed and included in the PATH")
	}
}

// if file point to a Dir
func UpdateAndSync(fileName string) {

	GitPkg := git.GitDir{fileName, false}
	if e := GitPkg.Pull(); e != nil {
		fmt.Fprintf(writer, "%s", e)
	}
	if !GitPkg.Updated && !ForceBuild {
		fmt.Fprintf(writer, "No updates found for package: %s\n", fileName)
		return
	}

	packages, err := pacman.MakePackage(fileName)
	if err != nil {
		log.Printf("%s", err)
	}

	for _, pkg := range packages {
		if e := pacman.UpdateRepo(RepoPath, pkg); e != nil {
			fmt.Fprintf(writer, "%s", e)
		}
	}
}

func main() {
	flag.Parse()

	var files []fs.DirEntry
	var e error

	if Filename != "" {
		fh, e := os.Stat(AurPath + Filename)
		if e != nil {
			log.Fatal("Unable to read File:", AurPath+Filename)
		}

		files = []fs.DirEntry{fs.FileInfoToDirEntry(fh)}
		ForceBuild = true
	} else {
		files, e = os.ReadDir(AurPath)
		if e != nil {
			log.Fatal("Unable to read Directory:", AurPath)
		}
	}

	for _, file := range files {
		UpdateAndSync(AurPath + file.Name())
	}
}
