package main

import (
	"bytes"
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

func Pull() {
	// client sends last mod times
	// check mod time
	// send list of updates
	// call serve file
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	project_file, err := ioutil.ReadFile(filepath.Join(pwd, ".disty", "project.json"))
	if err != nil {
		panic(err)
	}

	var project Project
	if err := json.Unmarshal(project_file, &project); err != nil {
		panic(err)
	}

	if project.Name == "" {
		panic("No project name\nPlease reinitialize the repo")
	}

	if project.Remote == "" {
		panic("No remote server\nPlease use `disty remote` to add remote server")
	}

	index, err := ioutil.ReadFile(filepath.Join(pwd, ".disty", "index.csv"))
	if err != nil {
		panic(err)
	}

	baseUrl := fmt.Sprintf("http://%s", project.Remote)
	url := fmt.Sprintf("%s/pull?project=%s", baseUrl, project.Name)

	client := http.Client{Timeout: time.Duration(5 * time.Second)}

	pullReq, err := http.NewRequest("GET", url, bytes.NewReader(index))
	if err != nil {
		panic(err)
	}

	indexResp, err := client.Do(pullReq)
	defer indexResp.Body.Close()
	if err != nil {
		panic(err)
	}

	indexBodyBytes, err := ioutil.ReadAll(indexResp.Body)
	if err != nil {
		panic(err)
	}

	indexBody := string(indexBodyBytes)
	files := strings.Split(indexBody, "|")

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		// DRY
		go func(url, project, file string) {
			defer wg.Done()
			resp, err := http.Get(fmt.Sprintf("%s/file?project=%s&filename=%s", url, project, file))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			respBodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			filenameIndex := strings.LastIndex(file, "\\")

			if filenameIndex == -1 {
				return
			}

			if _, err := os.Stat(filepath.Join(pwd, file[:filenameIndex])); os.IsNotExist(err) {
				err := os.MkdirAll(filepath.Join(pwd, file[:filenameIndex]), 0777)
				if err != nil {
					panic(err)
				}
			}

			f, err := os.Create(filepath.Join(pwd, file))
			defer f.Close()
			if err != nil {
				panic(err)
			}

			_, err = f.Write(respBodyBytes)
			if err != nil {
				panic(err)
			}
		}(baseUrl, project.Name, file)
	}

	wg.Wait()

	fmt.Println("[SUCCESS] Successfully pulled")
}
