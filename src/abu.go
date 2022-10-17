package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	TRAGET_REPO = "feature"
	MASTER_REPO = "master"
	DATE_LAYOUT = "20060102"
)

func main() {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	fp := createFilePointer()
	defer fp.Close()

	gitBranch(fp)
	gitPullBranch(nil, TRAGET_REPO, fp)
	gitPushBrunch(nil, TRAGET_REPO, fp)
	gitCheckOutBrunch(nil, MASTER_REPO, fp)
	gitBranch(fp)
	gitPullReleaseToTarget(nil, TRAGET_REPO, MASTER_REPO, fp)
	gitPushBrunch(nil, TRAGET_REPO, fp)
}

func createFilePointer() (fp *os.File) {
	day := time.Now()
	today_date := day.Format(DATE_LAYOUT)

	fileName := fmt.Sprintf("%s-%s.log", today_date, TRAGET_REPO)

	fp, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	return fp
}

func gitBranch(fp *os.File) {
	fp.WriteString("\n--- git branch ---\n")
	output, err := exec.Command("git", "branch").CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(errStr)
		os.Exit(1)
	}
	fp.WriteString(string(output))
}

func gitPullBranch(in io.Reader, repoName string, fp *os.File) {
	fp.WriteString("\n--- git pull --progress origin ---\n")
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
			errStr := fmt.Sprintf("%v\n", err)
			fp.WriteString(errStr)
			os.Exit(1)
		}
		fp.WriteString(string(output))
	default:
		fmt.Printf("do not git pull %s\n", repoName)
	}
}

func gitPushBrunch(in io.Reader, repoName string, fp *os.File) {
	refsRepo := fmt.Sprintf("refs/heads/%s:refs/heads/%s", repoName, repoName)
	fileWriteStr := fmt.Sprintf("\n--- git push --recurse-submodules=check origin %s ---\n", refsRepo)
	fp.WriteString(fileWriteStr)
	if in == nil {
		in = os.Stdin
	}
	fmt.Printf("git push %s? >", repoName)

	var inputValue string
	fmt.Fscan(in, &inputValue)
	switch inputValue {
	case "y", "Y", "yes", "YES":
		output, err := exec.Command("git", "push", "--recurse-submodules=check", "origin", refsRepo).CombinedOutput()
		if err != nil {
			errStr := fmt.Sprintf("%v\n", err)
			fp.WriteString(errStr)
			os.Exit(1)
		}
		fp.WriteString(string(output))
	default:
		fmt.Printf("do not git push %s\n", repoName)
	}
}

func gitCheckOutBrunch(in io.Reader, repoName string, fp *os.File) {
	fileWriteStr := fmt.Sprintf("\n--- git checkout %s ---\n", repoName)
	fp.WriteString(fileWriteStr)
	if in == nil {
		in = os.Stdin
	}
	fmt.Printf("git checkout %s? >", repoName)

	var inputValue string
	fmt.Fscan(in, &inputValue)
	switch inputValue {
	case "y", "Y", "yes", "YES":
		// git checkout
		output, err := exec.Command("git", "checkout", repoName).CombinedOutput()
		if err != nil {
			errStr := fmt.Sprintf("%v\n", err)
			fp.WriteString(errStr)
			os.Exit(1)
		}
		fp.WriteString(string(output))
	default:
		fmt.Printf("do not git checkout %s\n", repoName)
	}
}

func gitPullReleaseToTarget(in io.Reader, repoName string, masterRepoName string, fp *os.File) {
	refsMasterRepo := fmt.Sprintf("refs/heads/%s", masterRepoName)
	fileWriteStr := fmt.Sprintf("\n--- git pull --progress origin %s ---\n", refsMasterRepo)
	fp.WriteString(fileWriteStr)
	if in == nil {
		in = os.Stdin
	}
	fmt.Printf("git pull %s to %s? >", masterRepoName, repoName)

	var inputValue string
	fmt.Fscan(in, &inputValue)
	switch inputValue {
	case "y", "Y", "yes", "YES":
		output, err := exec.Command("git", "pull", "--progress", "origin", refsMasterRepo).CombinedOutput()
		if err != nil {
			errStr := fmt.Sprintf("%v\n", err)
			fp.WriteString(errStr)
			os.Exit(1)
		}
		fp.WriteString(string(output))
	default:
		fmt.Printf("do not git pull %s\n", masterRepoName)
	}

}
