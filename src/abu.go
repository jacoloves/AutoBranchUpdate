package main

import (
	"os"
	"path/filepath"
)

const (
	TRAGET_REPO = "test"
	MATER_REPO  = "master"
)

func main() {
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)
}
