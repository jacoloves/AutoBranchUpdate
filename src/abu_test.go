package main

import (
	"bytes"
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

type BranchInformation struct {
	Id             int      `json:"id"`
	MainRepository string   `json:"mainRepository"`
	LogRepository  string   `json:"logRepository"`
	MasterBranch   string   `json:"masterBranch"`
	RepositoryName string   `json:"repositoryName"`
	TargetBranches []string `json:"targetBranches"`
}

type BranchInformations struct {
	BranchInformations []BranchInformation `json:"settingArray"`
}

func Test_gitPullBranch(t *testing.T) {
	fmt.Println("--- Test gitPullBranch ---")
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
	fmt.Println("--- Test gitPushBranch ---")
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

func Test_gitCheckOutBrunch(T *testing.T) {
	fmt.Println("--- Test gitCheckOutBrunch ---")
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		in := bytes.NewBufferString("y")
		gitCheckOutBrunch(in, "feature")
	}

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		in := bytes.NewBufferString("n")
		gitCheckOutBrunch(in, "feature")
	}

	fmt.Println("--- input nil test ---")
	gitCheckOutBrunch(nil, "feature")

	fmt.Println("--- nil yes strings test ---")
	os.Chdir(prev)
	for _, input := range yesstrings {
		in := bytes.NewBufferString(input)
		gitCheckOutBrunch(in, "feature")
	}
}

func Test_gitPullReleaseToTarget(t *testing.T) {
	fmt.Println("--- Test gitPullReleaseToTarget ---")
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		in := bytes.NewBufferString("y")
		gitPullReleaseToTarget(in, "feature", "master")
	}

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		in := bytes.NewBufferString("n")
		gitPullReleaseToTarget(in, "feature", "master")
	}

	fmt.Println("--- input nil test ---")
	gitPullReleaseToTarget(nil, "feature", "master")

	fmt.Println("--- input yes strings test ---")
	os.Chdir(prev)
	for _, input := range yesstrings {
		in := bytes.NewBufferString(input)
		gitPullReleaseToTarget(in, "feature", "master")
	}

}

func Test_getConfigData(t *testing.T) {
	fmt.Println("--- Test getConfigData ---")

	fmt.Println("--- error ReadFile ---")
	_ = getConfigData("./notexist.json")

	fmt.Println("-- error Unmarshal ---")
	_ = getConfigData("./unmarshalng.json")

	fmt.Println("--- pass process ---")
	_ = getConfigData("./setting.json")
}
