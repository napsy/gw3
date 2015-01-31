package main

import (
	git "gitw3"
	"log"
)

func Job(project *Project, branchName string) error {
	var (
		logEntry []git.LogEntry
		gitHome  *git.GitRoot
		step     int = 1
		err      error
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
	gitHome = git.NewGit(project.Home)
	for ; step < 5; step++ {
		makeStep(step)
	}
	if err != nil {
		log.Printf("Failed at step %d: %v", step, err)
	}
	return nil
}
