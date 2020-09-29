package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Clone(url, dirpath, project string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/clone?project=%s", url, project))
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(dirpath, 0777) // make dir
	if err != nil {
		panic(err)
	}

	body := string(bodyBytes)
	files := strings.Split(body, "|")

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		// not sure about putting file io in goroutines
		go func(url, project, file string) {
			defer wg.Done()
			resp, err := http.Get(fmt.Sprintf("http://%s/file?project=%s&filename=%s", url, project, file))
			defer resp.Body.Close()
			if err != nil {
				panic(err)
			}

			respBodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			filenameIndex := strings.LastIndex(file, "\\")

			if _, err := os.Stat(filepath.Join(pwd, dirpath, file[:filenameIndex])); os.IsNotExist(err) {
				err := os.MkdirAll(filepath.Join(pwd, dirpath, file[:filenameIndex]), 0700)
				if err != nil {
					panic(err)
				}
			}

			f, err := os.Create(filepath.Join(pwd, dirpath, file))
			defer f.Close()
			if err != nil {
				panic(err)
			}

			_, err = f.Write(respBodyBytes)
			if err != nil {
				panic(err)
			}
		}(url, project, file)
	}
	wg.Wait()

	fmt.Println("Successfully cloned")
}
