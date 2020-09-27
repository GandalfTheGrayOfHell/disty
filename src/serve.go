package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func servePush(w http.ResponseWriter, r *http.Request, dir string) {
	query := r.URL.Query()
	fmt.Println(query)

	project, err := query["project"]
	if !err || len(project[0]) < 1 {
		w.WriteHeader(400) // Bad request
		return
	}

	// check if project exists on disk
	if _, err := os.Stat(filepath.Join(dir, project[0])); os.IsNotExist(err) {
		// project does not exist
		w.WriteHeader(500)
		return
	}

	file_name, err := query["filename"]
	filename := filepath.FromSlash(file_name[0])

	if !err || len(file_name[0]) < 1 {
		w.WriteHeader(400) // Bad request
		return
	}

	lmod, err := query["modtime"] // local mod time
	if !err || len(lmod) < 1 {
		w.WriteHeader(400)
	}

	// check if file exists
	file_info, err1 := os.Stat(filepath.Join(dir, project[0], filename))

	if os.IsNotExist(err1) || file_info.IsDir() { // file does not exist or is a directory
		if req_to_file(r, filepath.Join(dir, project[0], filename), 0777) != nil { // perform write on file here
			w.WriteHeader(500)
			return
		}

		err := add_index_file_mod(filepath.Join(dir, project[0], ".disty", "index.csv"), filename, lmod[0])
		if err != nil {
			w.WriteHeader(500)
			return
		}

	} else if !os.IsNotExist(err1) && !file_info.IsDir() { // file exists and not a directory
		rmod, err := get_index_file_modtime(filepath.Join(dir, project[0], ".disty", "index.csv"), filename) // remote mod time
		if err != nil {
			w.WriteHeader(500)
			return
		}

		irmod, err := strconv.Atoi(rmod)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		ilmod, err := strconv.Atoi(lmod[0])
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if irmod < ilmod {
			if req_to_file(r, filepath.Join(dir, project[0], filename), 0777) != nil { // perform write on file
				w.WriteHeader(500)
				return
			}

			// update mod time of file in index
			err := update_index_file_mod(filepath.Join(dir, project[0], ".disty", "index.csv"), filename, lmod[0], "0")
			if err != nil {
				w.WriteHeader(500)
				return
			}

		}
	} else if err1 != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func servePull(w http.ResponseWriter, r *http.Request, dir string) {
	query := r.URL.Query()

	project, err := query["project"]
	if !err || len(project[0]) < 1 {
		w.WriteHeader(400) // Bad request
		return
	}

	// check if project exists on disk
	if _, err := os.Stat(filepath.Join(dir, project[0])); os.IsNotExist(err) {
		// project does not exist
		w.WriteHeader(500)
		return
	}

	// read client index
	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		w.WriteHeader(400)
		return
	}

	req_reader := csv.NewReader(bytes.NewReader(body))

	req_records, err1 := req_reader.ReadAll()
	if err1 != nil {
		w.WriteHeader(400)
		return
	}

	req_map := make(map[string]string)
	for _, record := range req_records {
		req_map[record[0]] = record[1]
	}

	// read local index file for a project
	index, err1 := ioutil.ReadFile(filepath.Join(dir, project[0], "index.csv"))
	if err1 != nil {
		w.WriteHeader(500)
		return
	}

	index_reader := csv.NewReader(bytes.NewReader(index))
	index_records, err1 := index_reader.ReadAll()
	if err1 != nil {
		w.WriteHeader(500)
		return
	}

	var updation []string

	for _, record := range index_records {
		value, prs := req_map[record[0]]

		if prs == true {
			file_mod, err := strconv.Atoi(value)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			index_mod, err := strconv.Atoi(record[1])
			if err != nil {
				w.WriteHeader(500)
				return
			}

			// check for changes with repo
			if index_mod > file_mod {
				updation = append(updation, record[0])
			}
		} else {
			w.WriteHeader(404)
			return
		}
	}

	// send updation file

	w.Write([]byte(strings.Join(updation[:], "|")))
}

func serveFile(w http.ResponseWriter, r *http.Request, dir string) {
	query := r.URL.Query()

	project, err := query["project"]
	if !err || len(project[0]) < 1 {
		w.WriteHeader(400) // Bad request
		return
	}

	file, err := query["file"]
	if !err || len(file[0]) < 1 {
		w.WriteHeader(400)
		return
	}

	// check if file exist
	if _, err := os.Stat(filepath.Join(dir, project[0], file[0])); os.IsNotExist(err) {
		// project does not exist
		w.WriteHeader(500)
		return
	}

	bytes, err1 := ioutil.ReadFile(filepath.Join(dir, project[0], file[0]))
	if err1 != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(bytes)
}

func serveClone(w http.ResponseWriter, r *http.Request, dir string) {
	pwd, err := os.Getwd()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	query := r.URL.Query()

	project, err1 := query["project"]
	if !err1 || len(project[0]) < 1 {
		w.WriteHeader(400) // Bad request
		return
	}

	var files []string

	err = filepath.Walk(filepath.Join(dir, project[0]), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			files = append(files, strings.Replace(path, filepath.Join(pwd, dir), "", 1))
		}

		return nil
	})

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write([]byte(strings.Join(files[:], "|")))
}

func Serve(port int, dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// dir does not exist
		if os.MkdirAll(dir, 0777) != nil {
			panic("[ERROR] Could not create dir: " + dir)
		}
	}

	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		servePush(w, r, dir)
	})
	http.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		servePull(w, r, dir)
	})
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, r, dir)
	})
	http.HandleFunc("/clone", func(w http.ResponseWriter, r *http.Request) {
		serveClone(w, r, dir)
	})

	if http.ListenAndServe(":"+strconv.Itoa(port), nil) != nil {
		panic("[ERROR] Could not start server on port " + string(port))
	}
}
