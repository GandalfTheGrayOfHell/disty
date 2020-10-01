package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Project struct {
	Remote string `json:"remote"`
	Name   string `json:"name"`
}

// Writes a Request Body to a file
func req_to_file(r *http.Request, file string, perm os.FileMode) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, body, perm)
	if err != nil {
		return err
	}

	return nil
}

// gets the mod time of a file in the index
func get_index_file_modtime(index string, filename string) (string, error) {
	reader, err := ioutil.ReadFile(index)
	if err != nil {
		return "", err
	}

	r := csv.NewReader(bytes.NewReader(reader))

	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	var rec string = "-1"

	// record = [path, modtime]
	for _, record := range records {
		if record[0] == filename {
			rec = record[1]
		}
	}

	return rec, nil
}

// updates the mod time for a TRACKED file in an index
func update_index_file_mod(index string, filename string, modtime string, update string) error {
	reader, err := ioutil.ReadFile(index)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(reader))

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	for i := range records {
		if records[i][0] == filename {
			records[i][1] = modtime
			records[i][2] = update
			break
		}
	}

	f, err := os.OpenFile(index, os.O_SYNC|os.O_WRONLY, 0777)
	defer f.Close()
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll(records)

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func reset_index_updates(index string) error {
	reader, err := ioutil.ReadFile(index)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(reader))

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	for i := range records {
		records[i][2] = "0"
	}

	f, err := os.OpenFile(index, os.O_SYNC|os.O_WRONLY, 0777)
	defer f.Close()
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll(records)

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

// adds an UNTRACKED file with its mod time to index
func add_index_file_mod(index string, filename string, modtime string) error {
	reader, err := ioutil.ReadFile(index)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(reader))

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	records = append(records, []string{filename, modtime, "1"})

	f, err := os.OpenFile(index, os.O_SYNC|os.O_WRONLY, 0777)
	defer f.Close()
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll(records)

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func check_project_modified(project string) bool {
	reader, err := ioutil.ReadFile(filepath.Join(project, ".disty", "index.csv"))
	if err != nil {
		panic("[ERROR] Could not read index")
	}

	r := csv.NewReader(bytes.NewReader(reader))

	records, err := r.ReadAll()
	if err != nil {
		panic("[ERROR] Could not parse CSV records")
	}

	filemap := make(map[string]string)

	for _, record := range records {
		filemap[record[0]] = record[1]
	}

	var modified []string
	var untracked []string

	filepath.Walk(project, func(path string, info os.FileInfo, err error) error {
		rel_path := strings.Replace(path, project, "", 1)

		if err != nil {
			panic("[ERROR] `Walk` error for file: " + rel_path)
		}

		if info.IsDir() && (info.Name() == ".git" || info.Name() == ".disty") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			file_info, err := os.Stat(path)
			if err != nil {
				panic("[ERROR] Could not find specified file: " + rel_path)
			}

			value, prs := filemap[rel_path]

			if prs == false {
				untracked = append(untracked, rel_path)
			} else {
				file_mod := int(file_info.ModTime().Unix())

				index_mod, err := strconv.Atoi(value)
				if err != nil {
					panic("[ERROR] Could not convert indexed Mod time to int")
				}

				if file_mod > index_mod {
					modified = append(modified, rel_path)
				}
			}
		}

		return nil
	})

	if len(modified) == 0 && len(untracked) == 0 {
		return false
	} else {
		return true
	}
}
