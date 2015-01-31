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
type PushStruct struct {
	Ref        string `json:",after"`
	Repository struct {
		Name        string `json:",name"`
		URL         string `json:",url"`
		Description string `json:",description"`
		Homepage    string `json:",homepage"`
	} `json:",repository"`
}

type Project struct {
	Home string
}

type Config struct {
	Host     string
	Port     string
	Projects map[string]Project
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
	http.HandleFunc("/incoming", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Couldn't read body: %v", err)
			return
		}
		defer r.Body.Close()

		push := PushStruct{}
		if err := json.Unmarshal(body, &push); err != nil {
			log.Printf("Couldn't read body: %v", err)
			return
		}
		project, ok := GlobalConfig.Projects[push.Repository.Name]
		if !ok {
			http.Error(w, ErrProject.Error(), 500)
			return
		}
		branchName := push.Ref[strings.LastIndex(push.Ref, "/"):]
		if err = Job(&project, branchName); err != nil {
			http.Error(w, ErrProject.Error(), 500)
		}
	})
	if err := LoadConfig(); err != nil {
		log.Fatal("Unable to read configuraton: %v", err)
	}
	log.Printf("Listening on %s:%s\n", GlobalConfig.Host, GlobalConfig.Port)
	log.Fatal(http.ListenAndServe(GlobalConfig.Host+":"+GlobalConfig.Port, nil))
}

func init() {
}
