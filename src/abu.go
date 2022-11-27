package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

const (
	TRAGET_REPO = "feature"
	MASTER_REPO = "master"
	DATE_LAYOUT = "20060102"
	CONFIG_FILE = "./setting.json"
)

type BranchInformation struct {
	Id             int      `json:"id"`
	MainRepository string   `json:"mainRepository"`
	LogRepository  string   `json:"logRepository"`
	MasterBranch   string   `json:"masterBranch"`
	RepositoryName string   `json:"repositoryName"`
	TargetBranches []string `json:"targetBranches"`
}

type BranchInformationArray struct {
	BranchInformationArray []BranchInformation `json:"branchInformationArray"`
}

func main() {
	// current direcotry get
	prev, err := filepath.Abs(".")
	if err != nil {
		os.Exit(1)
	}
	defer os.Chdir(prev)

	// Json file data get
	configArray := getConfigData(CONFIG_FILE)

	for _, data := range configArray.BranchInformationArray {
		masterBranchOperationFlg := true
		fmt.Printf("======= %s repository's branches update! ======\n", data.RepositoryName)
		processErrFlg := createLogDir(data.LogRepository)
		if processErrFlg {
			fmt.Println("createLogDir func failed")
			continue
		}

		// change directory
		os.Chdir(data.MainRepository)

		for _, branch := range data.TargetBranches {
			fmt.Printf("%s ... ", branch)
			// filepointer creates
			fp, processErrFlg := createFilePointer(data.LogRepository, branch)
			if processErrFlg {
				printResultColor(processErrFlg)
				continue
			}
			defer fp.Close()

			// For the first time, the operation of masterBranch is executed.
			if masterBranchOperationFlg {
				// master branch checkout
				processErrFlg = gitCheckOutBranch(data.MasterBranch, fp)
				if processErrFlg {
					printResultColor(processErrFlg)
					continue
				}
				// master branch pull
				processErrFlg = gitPullBranch(data.MainRepository, fp)
				if processErrFlg {
					printResultColor(processErrFlg)
					continue
				}

				// flg data chnage
				masterBranchOperationFlg = false
			}

			// target branch checkout
			processErrFlg = gitCheckOutBranch(branch, fp)
			if processErrFlg {
				printResultColor(processErrFlg)
				continue
			}

			// target branch pull
			processErrFlg = gitPullBranch(branch, fp)
			if processErrFlg {
				printResultColor(processErrFlg)
				continue
			}

			// target branch push
			processErrFlg = gitPushBranch(branch, fp)
			if processErrFlg {
				printResultColor(processErrFlg)
				continue
			}

			// git pull master branch to target branch
			gitPullReleaseToTarget(data.MasterBranch, fp)
			if processErrFlg {
				printResultColor(processErrFlg)
				continue
			}

			// target branch push
			processErrFlg = gitPushBranch(branch, fp)
			if processErrFlg {
				printResultColor(processErrFlg)
				continue
			}

			printResultColor(processErrFlg)

		}
	}
	fmt.Println("!!!AutoBranchUpdate complete!!!")
}

func printResultColor(errFlg bool) {
	if errFlg {
		fmt.Print("\x1b[38;2;255;0;0m")
		fmt.Println("NG!!")
		fmt.Print("\x1b[0m")
	} else {
		fmt.Print("\x1b[38;2;0;126;0m")
		fmt.Println("OK!!")
		fmt.Print("\x1b[0m")
	}
}

func replaceTildeToHomedir(dirName string) string {
	usr, _ := user.Current()
	replacedDirName := strings.Replace(dirName, "~", usr.HomeDir, 1)

	return replacedDirName
}

func createLogDir(createLogDir string) bool {
	day := time.Now()
	today_date := day.Format(DATE_LAYOUT)

	createLogDir = replaceTildeToHomedir(createLogDir)

	os.Chdir(createLogDir)
	if err := os.Mkdir(today_date, 0777); err != nil {
		fmt.Println(err)
		return true
	}

	return false
}

func getConfigData(configFileName string) BranchInformationArray {
	raw, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var branchDatas BranchInformationArray
	if err = json.Unmarshal(raw, &branchDatas); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return branchDatas
}

func createFilePointer(logDirName string, branchName string) (fp *os.File, errFlg bool) {
	day := time.Now()
	today_date := day.Format(DATE_LAYOUT)

	logDirName = replaceTildeToHomedir(logDirName)

	fileName := fmt.Sprintf("%s/%s/%s.log", logDirName, today_date, branchName)

	fp, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil, true
	}

	return fp, false
}

func gitPullBranch(repoName string, fp *os.File) (errFlg bool) {
	fp.WriteString("\n--- git pull --progress origin ---\n")

	output, err := exec.Command("git", "pull", "--progress", "origin").CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(string(output))
		fp.WriteString(errStr)
		return true
	}
	fp.WriteString(string(output))
	return false
}

func gitPushBranch(repoName string, fp *os.File) (errFlg bool) {
	refsRepo := fmt.Sprintf("refs/heads/%s:refs/heads/%s", repoName, repoName)
	fileWriteStr := fmt.Sprintf("\n--- git push --recurse-submodules=check origin %s ---\n", refsRepo)
	fp.WriteString(fileWriteStr)

	output, err := exec.Command("git", "push", "--recurse-submodules=check", "origin", refsRepo).CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(string(output))
		fp.WriteString(errStr)
		return true
	}
	fp.WriteString(string(output))
	return false
}

func gitCheckOutBranch(repoName string, fp *os.File) (errFlg bool) {
	fileWriteStr := fmt.Sprintf("\n--- git checkout %s ---\n", repoName)
	fp.WriteString(fileWriteStr)

	// git checkout
	output, err := exec.Command("git", "checkout", repoName).CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(string(output))
		fp.WriteString(errStr)
		return true
	}
	fp.WriteString(string(output))
	return false
}

func gitPullReleaseToTarget(masterRepoName string, fp *os.File) (errFlg bool) {
	refsMasterRepo := fmt.Sprintf("+refs/heads/%s", masterRepoName)
	fileWriteStr := fmt.Sprintf("\n--- git pull --progress origin %s ---\n", masterRepoName)
	fp.WriteString(fileWriteStr)

	output, err := exec.Command("git", "pull", "--progress", "origin", refsMasterRepo).CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(string(output))
		fp.WriteString(errStr)
		return true
	}
	fp.WriteString(string(output))
	return false
}
