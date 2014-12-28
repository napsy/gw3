package main

import (
	"fmt"
	"gitw3"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	projects := map[string]*gitw3.GitRoot{"koala2": gitw3.NewGit("/home/luka/work/koala2")}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, projects)
	})
	r.HandleFunc("/{repo}/", func(w http.ResponseWriter, r *http.Request) {
		route := mux.Vars(r)
		name := route["repo"]

		r.ParseForm()
		branch := r.FormValue("branch")

		git, found := projects[name]
		if !found {
			fmt.Fprintf(w, "repository not found")
			return
		}
		if err := git.Checkout(branch); err != nil {
			log.Println(err)
		}
		logs, err := git.Log(30)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		branches, err := git.Branches()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		t, _ := template.ParseFiles("templates/logs.html")
		repository := struct {
			Name     string
			Branch   string
			Branches []string
			Logs     []gitw3.LogEntry
		}{Name: name, Branch: branch, Branches: branches, Logs: logs}
		t.Execute(w, repository)
	})
	r.HandleFunc("/{repo}/{commit}", func(w http.ResponseWriter, r *http.Request) {
		route := mux.Vars(r)
		name := route["repo"]
		commitId := route["commit"]

		git, found := projects[name]
		if !found {
			fmt.Fprintf(w, "repository not found")
			return
		}
		commit, err := git.GetCommit(commitId)
		if err != nil {
			log.Println(err)
		}
		t, _ := template.ParseFiles("templates/commit.html")
		t.Execute(w, commit)

	})
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
