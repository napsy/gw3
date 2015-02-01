package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
)

// GitLab PUSH structure
type GitlabPush struct {
	After      string `json:",after"`
	Ref        string `json:",ref"`
	Repository struct {
		Name        string `json:",name"`
		URL         string `json:",url"`
		Description string `json:",description"`
		Homepage    string `json:",homepage"`
	} `json:",repository"`
}

type Config struct {
	Host     string
	Port     string
	LogsDir  string
	Projects Projects
}

var (
	ErrProject = errors.New("project not found")
)

var (
	ConfigPath   = flag.String("c", "config.toml", "configuration file")
	GlobalConfig Config
)

func LoadConfig() error {
	if _, err := toml.DecodeFile(*ConfigPath, &GlobalConfig); err != nil {
		return err
	}
	return nil
}
func main() {
	http.HandleFunc("/gitlab/incoming", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Couldn't read body: %v", err)
			return
		}
		defer r.Body.Close()

		push := GitlabPush{}
		if err := json.Unmarshal(body, &push); err != nil {
			log.Printf("Couldn't read body: %v", err)
			return
		}
		// TODO(napsy): generalize push structure
		projectName := LoadedProjects.Get(push.Repository.Name)
		if projectName == nil {
			http.Error(w, ErrProject.Error(), 500)
			return
		}
		branchName := push.Ref[strings.LastIndex(push.Ref, "/"):]
		commitID := push.After
		job := NewJob(projectName, map[string]string{"branch": branchName, "commit": commitID})
		if IncomingJobs.Push(job) == false {
			log.Printf("Job queue full, job request for %s:%s was not enqueued", projectName, commitID)
		}
	})
	if err := LoadConfig(); err != nil {
		log.Fatalf("Unable to read configuraton: %v", err)
	}
	LoadedProjects = GlobalConfig.Projects
	log.Printf("Listening on %s:%s\n", GlobalConfig.Host, GlobalConfig.Port)
	log.Fatal(http.ListenAndServe(GlobalConfig.Host+":"+GlobalConfig.Port, nil))
}

func init() {
}
