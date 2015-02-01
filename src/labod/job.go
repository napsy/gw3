package main

import (
	"fmt"
	git "gitw3"
	"io"
	"log"
	"os"
	"time"
)

type Job struct {
	project       *Project
	args          map[string]string
	StartTime     time.Time
	statusUpdates bool // If enabled, status changes will be pushed
	statusChan    chan string
	outputWriter  io.Writer
}

func NewJob(project *Project, args map[string]string) *Job {
	return &Job{project: project, args: args}
}

func (job *Job) SyncWithRemote() (*git.LogEntry, error) {
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
			err = gitHome.Pull(job.args["branch"])
		case 2:
			err = gitHome.Checkout(job.args["commit"])
		case 3:
			logEntry, err = gitHome.Log(1)
		}
		step++
	}
	gitHome = git.NewGit(job.project.Home)
	for ; step < 4; step++ {
		makeStep(step)
	}
	if err != nil {
		return nil, fmt.Errorf("step %d: %v", step, err)
	}
	return &logEntry[0], nil
}
func (job *Job) Worker() {
	var err error
	outputFilename := fmt.Sprintf("%s/%s-%s.txt", GlobalConfig.LogsDir, job.StartTime.Format(time.RFC3339), job.project.name)
	job.outputWriter, err = os.Create(outputFilename)
	if err != nil {
		log.Printf("Couldn't create output file: %v")
		log.Printf("Dumping output to stdout ...")
		job.outputWriter = os.Stdout
	}
	// First we need to sync our local repository with the changes
	commit, err := job.SyncWithRemote()
	if err != nil {
		log.Printf("Failed syncing with remote: %v", err)
		return
	}
	log.Printf("worker: working on %s: '%s' by %s ...", commit.Id, commit.Message, commit.Author)
}
func (job *Job) Start() {
	job.StartTime = time.Now()
	go job.Worker()
}
