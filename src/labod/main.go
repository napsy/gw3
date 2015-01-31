package main

import (
	"errors"
	"fmt"
	git "gitw3"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

type Project struct {
	Home string
}

type Labod struct {
	Projects map[string]Project
}

var (
	ErrProject = errors.New("project not found")
)

var (
	GlobalConfig Labod
)

func LoadConfig() {
	if _, err := toml.DecodeFile("config.toml", &GlobalConfig); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v\n", GlobalConfig)
}
func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		var (
			projectName string
			branchName  string
			project     Project
			logEntry    []git.LogEntry
			gitHome     *git.GitRoot
			step        int = 1
			ok          bool
			err         error
		)
		makeStep := func(step int) {
			if err != nil {
				return
			}
			switch step {
			case 1:
				err = gitHome.Pull(branchName)
			case 2:
				err = gitHome.Checkout(branchName)
			case 3:
				err = gitHome.Pull(branchName)
			case 4:
				logEntry, err = gitHome.Log(1)
			}
			step++
		}
		if project, ok = GlobalConfig.Projects[projectName]; ok {
			http.Error(w, ErrProject.Error(), 500)
			return
		}
		gitHome = git.NewGit(project.Home)
		for ; step < 5; step++ {
			makeStep(step)
		}
		if err != nil {
			log.Printf("Failed at step %d: %v", step, err)
		}
	})
	LoadConfig()
	log.Fatal(http.ListenAndServe(":8084", nil))
}

func init() {
}
