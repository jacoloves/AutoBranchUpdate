package main

import (
	"fmt"
	"io"
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
	gitPullBranch(nil, TRAGET_REPO)
	gitPushBrunch(nil, TRAGET_REPO)
}

func gitBranch() {
	output, err := exec.Command("git", "branch").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	fmt.Printf("%s", output)
}

func gitPullBranch(in io.Reader, repoName string) {
	if in == nil {
		in = os.Stdin
	}
	fmt.Printf("git pull %s? >", repoName)

	var inputValue string
	fmt.Fscan(in, &inputValue)
	switch inputValue {
	case "y", "Y", "yes", "YES":
		output, err := exec.Command("git", "pull", "--progress", "origin").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		fmt.Printf("%s", output)
	default:
		fmt.Printf("do not git pull %s\n", repoName)
	}
}

func gitPushBrunch(in io.Reader, repoName string) {
	if in == nil {
		in = os.Stdin
	}
	fmt.Printf("git push %s? >", repoName)

	refsRepo := fmt.Sprintf("refs/heads/%s:refs/heads/%s", repoName, repoName)
	var inputValue string
	fmt.Fscan(in, &inputValue)
	switch inputValue {
	case "y", "Y", "yes", "YES":
		output, err := exec.Command("git", "push", "--recurse-submodules=check", "origin", refsRepo).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		fmt.Printf("%s", output)
	default:
		fmt.Printf("do not git push %s\n", repoName)
	}
}
