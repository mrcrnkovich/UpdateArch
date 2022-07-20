package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

var filename string

func init() {
	flag.StringVar(&filename, "filename", "", "file path to open")
	flag.Parse()
}

func sync_dir(work_dir, file_name string, ch chan<- string) {

	fmt.Printf("Now syncing: %s\n", file_name)

	git_pack := work_dir + file_name

	// Sync git files
	output, _ := exec.Command("git", "-C", git_pack, "pull").Output()
	fmt.Printf("%s", output)

	notNew, _ := regexp.Match(`up to date`, []byte(output))

	output, _ = exec.Command("git", "-C", git_pack, "pull").Output()

	os.Chdir(git_pack)

	// make new package
	output, _ = exec.Command("makepkg", "--noconfirm", "--cleanbuild", "--syncdeps", "--rmdeps", "--clean", "--force").Output()

	// add to repo
	output, _ = exec.Command("repo-add", "--quiet", "--new", "--remove").Output()

	//clean up repo and dirs?

	if notNew {
		ch <- fmt.Sprintf("%s", "completed")
	} else {
		ch <- fmt.Sprintf("%s", "new version found")
	}
}

func main() {

	AUR_DIR := "/home/mike/build/aur/"

	ch := make(chan string)
	files, _ := os.ReadDir(AUR_DIR)
	for _, file := range files {
		go sync_dir(AUR_DIR, file.Name(), ch)
	}

	for range files {
		fmt.Println(<-ch)
	}
}
