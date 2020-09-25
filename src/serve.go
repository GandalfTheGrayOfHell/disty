package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func push(w http.ResponseWriter, r *http.Request, dir string) {
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

	track, err := query["track"]
	if !err || len(track) < 1 {
		w.WriteHeader(400)
	}

	lmod, err := query["modtime"] // local mod time
	if !err || len(lmod) < 1 {
		w.WriteHeader(400)
	}

	// check if file exists
	file_info, err1 := os.Stat(filepath.Join(dir, project[0], filename))

	if os.IsNotExist(err1) || file_info.IsDir() { // file does not exist or is a directory
		if track[0] == "TRACKED" { // file is tracked but does not exist on the disk
			w.WriteHeader(500)
			return
		} else if track[0] == "UNTRACKED" { // file is new to the repo
			if req_to_file(r, filepath.Join(dir, project[0], filename), 0777) != nil { // perform write on file here
				w.WriteHeader(500)
				return
			}

			err := add_index_file_mod(filepath.Join(dir, project[0], ".disty", "index.csv"), filename, lmod[0])
			if err != nil {
				w.WriteHeader(500)
				return
			}
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
			err := update_index_file_mod(filepath.Join(dir, project[0], ".disty", "index.csv"), filename, lmod[0])
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

func pull(w http.ResponseWriter, r *http.Request, dir string) {
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
	// check for changes with repo
	// send updation file
	// client sends individual requests
}

func clone(w http.ResponseWriter, r *http.Request) {

}

func Serve(port int, dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// dir does not exist
		if os.MkdirAll(dir, 0777) != nil {
			panic("[ERROR] Could not create dir: " + dir)
		}
	}

	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		push(w, r, dir)
	})
	http.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		pull(w, r, dir)
	})
	http.HandleFunc("/clone", clone)

	if http.ListenAndServe(":"+strconv.Itoa(port), nil) != nil {
		panic("[ERROR] Could not start server on port " + string(port))
	}
}
