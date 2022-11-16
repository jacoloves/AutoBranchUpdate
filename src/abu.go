package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	TRAGET_REPO = "feature"
	MASTER_REPO = "master"
	DATE_LAYOUT = "20060102"
	CONFIG_FILE = "./setting.json"
)

type SettingData struct {
	Id             int      `json:"id"`
	MainRepository string   `json:"mainRepository"`
	LogRepository  string   `json:"logRepository"`
	MasterBranch   string   `json:"masterBranch"`
	RepositoryName string   `json:"repositoryName"`
	TargetBranches []string `json:"targetBranches"`
}

type SettingArray struct {
	SettingArray []SettingData `json:"settingArray"`
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

	for _, data := range configArray.SettingArray {
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
				fmt.Println("NG!!")
				continue
			}
			defer fp.Close()

			// For the first time, the operation of masterBranch is executed.
			if masterBranchOperationFlg {
				// master branch checkout
				processErrFlg = gitCheckOutBranch(data.MasterBranch, fp)
				if processErrFlg {
					fmt.Println("NG!!")
					continue
				}
				// master branch pull
				processErrFlg = gitPullBranch(data.MainRepository, fp)
				if processErrFlg {
					fmt.Println("NG!!")
					continue
				}

				// flg data chnage
				masterBranchOperationFlg = false
			}

			// target branch checkout
			processErrFlg = gitCheckOutBranch(branch, fp)
			if processErrFlg {
				fmt.Println("NG!!")
				continue
			}

			// target branch pull
			processErrFlg = gitPullBranch(branch, fp)
			if processErrFlg {
				fmt.Println("NG!!")
				continue
			}

			// target branch push
			processErrFlg = gitPushBranch(branch, fp)
			if processErrFlg {
				fmt.Println("NG!!")
				continue
			}

			// git pull master branch to target branch
			gitPullReleaseToTarget(data.MasterBranch, fp)
			if processErrFlg {
				fmt.Println("NG!!")
				continue
			}

			// target branch push
			processErrFlg = gitPushBranch(branch, fp)
			if processErrFlg {
				fmt.Println("NG!!")
				continue
			}

			fmt.Println("Ok!!")

		}
	}
	fmt.Println("!!!AutoBranchUpdate complete!!!")
}

func createLogDir(createLogDir string) bool {
	day := time.Now()
	today_date := day.Format(DATE_LAYOUT)

	os.Chdir(createLogDir)
	if err := os.Mkdir(today_date, 0777); err != nil {
		fmt.Println(err)
		return true
	}

	return false
}

func getConfigData(configFileName string) SettingArray {
	raw, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var settingDatas SettingArray
	if err = json.Unmarshal(raw, &settingDatas); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return settingDatas
}

func createFilePointer(logDirName string, branchName string) (fp *os.File, errFlg bool) {
	day := time.Now()
	today_date := day.Format(DATE_LAYOUT)

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
		fp.WriteString(errStr)
		return true
	}
	fp.WriteString(string(output))
	return false
}
