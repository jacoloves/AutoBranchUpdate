package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var testdirs = []string{
	".",                        // current direcotry git branch exist
	"/home",                    // home direcotry git branch not exist
	"/home/tests/test_abu_dir", // test direcotry git branch exist but remote repository no exist
}

type TestBranchInformation struct {
	Id             int      `json:"id"`
	MainRepository string   `json:"mainRepository"`
	LogRepository  string   `json:"logRepository"`
	MasterBranch   string   `json:"masterBranch"`
	RepositoryName string   `json:"repositoryName"`
	TargetBranches []string `json:"targetBranches"`
}

type TestBranchInformationArray struct {
	TestBranchInformationArray []TestBranchInformation `json:"branchInformationArray"`
}

func Test_gitPullBranch(t *testing.T) {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	fp, err := os.Create("test_abu.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer fp.Close()

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)
		_ = gitPullBranch("feature", fp)
	}

}

func Test_gitPushBranch(t *testing.T) {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	fp, err := os.Create("test_abu.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer fp.Close()

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		_ = gitPushBranch("feature", fp)
	}
}

func Test_gitCheckOutBranch(T *testing.T) {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	fp, err := os.Create("test_abu.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer fp.Close()

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		_ = gitCheckOutBranch("feature", fp)
	}

}

func Test_gitPullReleaseToTarget(t *testing.T) {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	fp, err := os.Create("test_abu.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer fp.Close()

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		_ = gitPullReleaseToTarget("feature", fp)
	}

}

func Test_getConfigData(t *testing.T) {
	_, _ = getConfigData("./notexist.json")

	_, _ = getConfigData("./unmarshalng.json")

	_, _ = getConfigData("./setting.json")

}

func Test_printResultColor(t *testing.T) {
	printResultColor(true)
	printResultColor(false)
}

func Test_replaceTildeToHomedir(t *testing.T) {
	_ = replaceTildeToHomedir("~/test/dir/tilde")
	_ = replaceTildeToHomedir("/home/test/dir/notilde")
}