package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	TRAGET_REPO = "feature"
	MATER_REPO  = "master"
)

func main() {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	gitBranch()
}

func gitBranch() {
	output, err := exec.Command("git", "branch").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	fmt.Printf("%s", output)
}
