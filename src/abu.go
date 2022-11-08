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

type Setting struct {
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
			fp := createFilePointer(data.LogRepository, branch)
			defer fp.Close()

			gitBranch(fp)
			gitPullBranch(TRAGET_REPO, fp)
			gitPushBranch(TRAGET_REPO, fp)
			gitCheckOutBranch(MASTER_REPO, fp)
			gitBranch(fp)
			gitPullReleaseToTarget(TRAGET_REPO, MASTER_REPO, fp)
			gitPushBranch(TRAGET_REPO, fp)

		}
		fmt.Println("Ok!!")
	}
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

func getConfigData(configFileName string) Setting {
	raw, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var settingDatas Setting
	if err = json.Unmarshal(raw, &settingDatas); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return settingDatas
}

func createFilePointer(logDirName string, branchName string) (fp *os.File) {
	day := time.Now()
	today_date := day.Format(DATE_LAYOUT)

	fileName := fmt.Sprintf("%s/%s/%s.log", logDirName, today_date, branchName)

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

func gitPullBranch(repoName string, fp *os.File) {
	fp.WriteString("\n--- git pull --progress origin ---\n")

	output, err := exec.Command("git", "pull", "--progress", "origin").CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(errStr)
		os.Exit(1)
	}
	fp.WriteString(string(output))
}

func gitPushBranch(repoName string, fp *os.File) {
	refsRepo := fmt.Sprintf("refs/heads/%s:refs/heads/%s", repoName, repoName)
	fileWriteStr := fmt.Sprintf("\n--- git push --recurse-submodules=check origin %s ---\n", refsRepo)
	fp.WriteString(fileWriteStr)

	output, err := exec.Command("git", "push", "--recurse-submodules=check", "origin", refsRepo).CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(errStr)
		os.Exit(1)
	}
	fp.WriteString(string(output))
}

func gitCheckOutBranch(repoName string, fp *os.File) {
	fileWriteStr := fmt.Sprintf("\n--- git checkout %s ---\n", repoName)
	fp.WriteString(fileWriteStr)

	// git checkout
	output, err := exec.Command("git", "checkout", repoName).CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(errStr)
		os.Exit(1)
	}
	fp.WriteString(string(output))
}

func gitPullReleaseToTarget(repoName string, masterRepoName string, fp *os.File) {
	refsMasterRepo := fmt.Sprintf("refs/heads/%s", masterRepoName)
	fileWriteStr := fmt.Sprintf("\n--- git pull --progress origin %s ---\n", refsMasterRepo)
	fp.WriteString(fileWriteStr)

	output, err := exec.Command("git", "pull", "--progress", "origin", refsMasterRepo).CombinedOutput()
	if err != nil {
		errStr := fmt.Sprintf("%v\n", err)
		fp.WriteString(errStr)
		os.Exit(1)
	}
	fp.WriteString(string(output))

}
