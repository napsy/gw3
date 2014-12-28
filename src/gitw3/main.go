package gitw3

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type GitRoot struct {
	root string
	exec func(rootDir string, args ...string) ([]string, error)
}

type Commit struct {
	Id      string
	Author  string
	Message string
	Diff    string
}

type LogEntry struct {
	Id      string
	Date    string
	Author  string
	Message string
}

var LocalGit = GitRoot{exec: gitExec}

func gitExec(rootDir string, args ...string) ([]string, error) {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "--no-pager")
	cmd.Args = append(cmd.Args, args...)
	cmd.Dir = rootDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		lines := strings.Split(string(out), "\n")
		return lines, err
	}
	lines := strings.Split(string(out), "\n")
	return lines, nil
}

// GetCommit gets the commit message and the diff for a commit
// with a particular id
func (git GitRoot) GetCommit(id string) ([]string, error) {
	msg, _ := git.exec(git.root, "show", id)
	return msg, nil
}

func (git GitRoot) Branches() ([]string, error) {
	msg, err := git.exec(git.root, "for-each-ref", "--sort=-committerdate", "refs/heads/", "--format='%(refname:short)'")
	if err != nil {
		return nil, fmt.Errorf("%v", msg)
	}
	for i := range msg {
		if len(msg[i]) == 0 {
			continue
		}
		msg[i] = msg[i][1 : len(msg[i])-1]
	}
	return msg, nil
}
func (git GitRoot) Checkout(branch string) error {
	msg, err := git.exec(git.root, "checkout", branch)
	if err != nil {
		return fmt.Errorf("%v", msg)
	}
	return nil

}
func (git GitRoot) Log(n int) ([]LogEntry, error) {
	nrLines := fmt.Sprintf("-n%d", n)
	// short id, full date, author, message
	msg, err := git.exec(git.root, "log", nrLines, "--date=relative", "--pretty=format:'%h,%ad,%aN,%s'")
	if err != nil {
		return nil, fmt.Errorf("%v", msg)
	}

	logEntries := []LogEntry{}
	for _, line := range msg {
		if len(line) < 1 {
			continue
		}
		tokens := strings.SplitN(line, ",", 4)
		if len(tokens) < 4 {
			return nil, errors.New("log format error")
		}
		if len(tokens[3]) > 80 {
			tokens[3] = fmt.Sprintf("%s ...", tokens[3][:80])
		}
		entry := LogEntry{Id: tokens[0][1:], Date: tokens[1], Author: tokens[2], Message: tokens[3][:len(tokens[3])-1]}
		logEntries = append(logEntries, entry)
	}
	return logEntries, nil
}

func NewGit(rootDir string) *GitRoot {
	return &GitRoot{exec: gitExec, root: rootDir}
}
