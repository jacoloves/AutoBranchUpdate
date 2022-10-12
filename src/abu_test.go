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

var yesstrings = []string{
	"y",
	"Y",
	"yes",
	"YES",
	"yas",
	"YeS",
	"yeS",
	"Yes",
	"year!",
	" y ",
	"y    ",
	"      y",
}

func Test_gitBranch(t *testing.T) {
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

		gitBranch()
	}
}

func Test_gitPullBranch(t *testing.T) {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	fmt.Println("--- directory change test y ---")
	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		in := bytes.NewBufferString("y")
		gitPullBranch(in, "feature")
	}

	fmt.Println("--- directory change test n ---")
	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		in := bytes.NewBufferString("n")
		gitPullBranch(in, "feature")
	}

	fmt.Println("--- input nil test ---")
	gitPullBranch(nil, "feature")

	fmt.Println("--- input yes strings test ---")
	os.Chdir(prev)
	for _, input := range yesstrings {
		in := bytes.NewBufferString(input)
		gitPullBranch(in, "feature")
	}
}

func Test_gitPushBranch(t *testing.T) {
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
		gitPushBrunch(in, "feature")
	}

	for _, dir := range testdirs {
		target, err := filepath.Abs(dir)
		if err != nil {
			os.Exit(1)
		}

		os.Chdir(target)

		in := bytes.NewBufferString("n")
		gitPushBrunch(in, "feature")
	}

	fmt.Println("--- input nil test ---")
	gitPushBrunch(nil, "feature")

	fmt.Println("--- input yes strings test ---")
	os.Chdir(prev)
	for _, input := range yesstrings {
		in := bytes.NewBufferString(input)
		gitPushBrunch(in, "feature")
	}
}

func Test_gitCheckOutBrunch(T *testing.T) {
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
