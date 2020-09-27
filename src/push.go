package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Push() {
	pwd, err := os.Getwd()
	if err != nil {
		panic("[ERROR] Could not get working directory")
	}

	if check_project_modified(pwd) == true {
		panic("There are changes in the repo\nPlease use `disty add` followed by appropiate files")
		return
	}

	projectIndex := strings.LastIndex(pwd, "\\")
	projectName := pwd[projectIndex+1:]

	fmt.Println(pwd, projectName)

	reader, err := ioutil.ReadFile(filepath.Join(pwd, ".disty", "index.csv"))
	if err != nil {
		panic("[ERROR] Could not read index")
	}

	r := csv.NewReader(bytes.NewReader(reader))

	records, err := r.ReadAll()
	if err != nil {
		panic("[ERROR] Could not read data from index")
	}

	project_file, err := ioutil.ReadFile(filepath.Join(pwd, ".disty", "project.json"))
	if err != nil {
		panic("[ERROR] Could not read project.json")
	}

	var project Project
	if json.Unmarshal(project_file, &project) != nil {
		panic("[ERROR] Invalid JSON during Unmarshal")
	}

	if project.Remote == "" {
		panic("There is no remote server\nUse `disty remote` to add a server")
	}

	chs := make(chan string, len(records))

	for _, record := range records {
		if record[2] == "1" {
			url := fmt.Sprintf("http://%s/push?project=%s&filename=%s&modtime=%s", project.Remote, projectName, record[0], record[1])

			body, err := ioutil.ReadFile(filepath.Join(pwd, record[0]))
			if err != nil {
				panic("[ERROR] Could not read file: " + filepath.Join(pwd, record[0]))
			}

			go make_request(url, "GET", body, chs)
		}
	}

	fmt.Println(<-chs)
	fmt.Println(<-chs)
	fmt.Println(<-chs)
	fmt.Println(<-chs)
}
