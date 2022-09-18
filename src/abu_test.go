package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitBranch(t *testing.T) {
	tests := [...]string{
		".",                        // current direcotry git branch exist
		"/home",                    // home direcotry git branch not exist
		"/home/tests/test_abu_dir", // test direcotry git branch exist but remote repository no exist
	}

	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	for _, dir := range tests {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		gitBranch()
	}
}
