package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Projects map[string]Project

type Project struct {
	Home    string
	PreCmd  string
	Cmd     string
	PostCmd string
	name    string
}

var (
	LoadedProjects Projects
)

func LoadProjects(projectsFile string) error {
	if _, err := toml.DecodeFile(projectsFile, &LoadedProjects); err != nil {
		return err
	}
	return nil
}

func (projects Projects) Get(name string) *Project {
	if project, ok := projects[name]; ok {
		return &project
	}
	return nil
}

func (projects Projects) List() []string {
	list := []string{}
	for name, project := range projects {
		info := fmt.Sprintf("%s: %+v", name, project)
		list = append(list, info)
	}
	return list
}
