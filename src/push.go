package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

	var wg sync.WaitGroup

	client := http.Client{Timeout: time.Duration(5 * time.Second)}

	for _, record := range records {
		if record[2] == "1" {
			body, err := ioutil.ReadFile(filepath.Join(pwd, record[0]))
			if err != nil {
				panic("[ERROR] Could not read file: " + filepath.Join(pwd, record[0]))
			}

			url := fmt.Sprintf("http://%s/push?project=%s&filename=%s&modtime=%s", project.Remote, projectName, record[0], record[1])

			wg.Add(1)
			go func(url string, method string, body []byte) {
				req, _ := http.NewRequest(method, url, bytes.NewReader(body))
				defer wg.Done()

				resp, err := client.Do(req)
				defer resp.Body.Close()
				if err != nil {
					panic("[ERROR] Could not make " + method + " request to " + url)
				}
			}(url, "GET", body)
		}
	}

	wg.Wait()

	if err := reset_index_updates(filepath.Join(pwd, ".disty", "index.csv")); err != nil {
		panic("[ERROR] Could not reset updates in index")
	}
}
