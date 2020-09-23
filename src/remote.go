package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO: remote verififcation

func Remote(remote string) {
	if remote == "" {
		panic("[ERROR] Remote URL cannot be empty")
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic("[ERROR] Could not get working directory path")
	}

	// TODO: check if project exists

	project_file, err := ioutil.ReadFile(filepath.Join(pwd, ".disty", "project.json"))
	if err != nil {
		panic("[ERROR] Could not read project.json")
	}

	var project Project
	if json.Unmarshal(project_file, &project) != nil {
		panic("[ERROR] Invalid JSON during Unmarshal")
	}

	project.Remote = remote
	project_json, err := json.MarshalIndent(project, "", "\t")
	if err != nil {
		panic("[ERROR] Could not marshal project")
	}

	if ioutil.WriteFile(filepath.Join(pwd, ".disty", "project.json"), project_json, 0777) != nil {
		panic("[ERROR] Could not write project.json")
	}

	fmt.Println("[SUCCESS] Added", remote, "as remote server")
}
