package main

import (
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
	Id               int      `json:"id"`
	MainRepository   string   `json:"mainRepository"`
	LogRepository    string   `json:"logRepository"`
	MasterBranch     string   `json:"masterBranch"`
	LogName          string   `json:"logName"`
	TargetRepository []string `json:"targetRepository"`
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

	fp := createFilePointer()
	defer fp.Close()

	gitBranch(fp)
	gitPullBranch(TRAGET_REPO, fp)
	gitPushBrunch(TRAGET_REPO, fp)
	gitCheckOutBrunch(MASTER_REPO, fp)
	gitBranch(fp)
	gitPullReleaseToTarget(TRAGET_REPO, MASTER_REPO, fp)
	gitPushBrunch(TRAGET_REPO, fp)
}

func getConfigData(configFileName string) {
	raw, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Println("--- config file load process failed ---")
		fmt.Println(err)
		os.Exit(1)
	}
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

func gitPushBrunch(repoName string, fp *os.File) {
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

func gitCheckOutBrunch(repoName string, fp *os.File) {
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
