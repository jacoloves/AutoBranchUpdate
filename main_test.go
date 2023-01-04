package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const TEST_DATE_LAYOUT = "20060102"

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

func Test_createLogDir(t *testing.T) {
	_ = createLogDir("/home/stanaka/released/AutoBranchUpdate/test")

	_ = createLogDir("/home/stanaka/released/AutoBranchUpdate/err_test/")

	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/20230104")
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/log/20230104")
}

func Test_createFilePointer(t *testing.T) {
	_ = createLogDir("/home/stanaka/released/AutoBranchUpdate/test")

	_, _ = createFilePointer("~/released/AutoBranchUpdate/test", "test_ok")

	_, _ = createFilePointer("~/released/AutoBranchUpdate/err_test", "test_ng")

	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/20230104")
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/log/20230104")
}

func Test_autoBranchUpdate(t *testing.T) {

	trueCase := BranchInformationArray{
		[]BranchInformation{
			{
				Id:             1,
				MainRepository: "~/released/AutoBranchUpdate",
				LogRepository:  "~/released/AutoBranchUpdate/test/log",
				MasterBranch:   "master",
				RepositoryName: "AutoBranchUpdate",
				TargetBranches: []string{"feature"},
			},
		},
	}

	errCreateLogDirCase := BranchInformationArray{
		[]BranchInformation{
			{
				Id:             2,
				MainRepository: "errCreateLogDirCase",
				LogRepository:  "",
				MasterBranch:   "errCreateLogDirCase",
				RepositoryName: "",
				TargetBranches: []string{"errCreateLogDirCase"},
			},
		},
	}

	errGitCheckOutBranchCase := BranchInformationArray{
		[]BranchInformation{
			{
				Id:             3,
				MainRepository: "~/released/AutoBranchUpdate",
				LogRepository:  "~/released/AutoBranchUpdate/test/log",
				MasterBranch:   "release",
				RepositoryName: "AutoBranchUpdate",
				TargetBranches: []string{""},
			},
		},
	}

	errGitPullBranchCase := BranchInformationArray{
		[]BranchInformation{
			{
				Id:             4,
				MainRepository: "home",
				LogRepository:  "~/released/AutoBranchUpdate/test/log",
				MasterBranch:   "master",
				RepositoryName: "AutoBranchUpdate",
				TargetBranches: []string{"test"},
			},
		},
	}

	_ = autoBranchUpdate(errCreateLogDirCase)
	_ = autoBranchUpdate(errGitCheckOutBranchCase)
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/20230104")
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/log/20230104")
	_ = autoBranchUpdate(errGitPullBranchCase)
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/20230104")
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/log/20230104")
	_ = autoBranchUpdate(trueCase)
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/20230104")
	_ = os.RemoveAll("/home/stanaka/released/AutoBranchUpdate/test/log/20230104")

}
