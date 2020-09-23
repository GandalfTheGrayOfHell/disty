package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"os"
)

type Project struct {
	Remote string `json:"remote"`
}

func (p *Project) Default() {
	p.Remote = ""
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
func get_file_modtime(index string, filename string) (string, error) {
	reader, err := ioutil.ReadFile(index)
	if err != nil {
		return "", err
	}

	r := csv.NewReader(bytes.NewReader(reader))

	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	var rec string = ""

	// record = [path, modtime]
	for _, record := range records {
		if record[0] == filename {
			rec = record[1]
		}
	}

	return rec, nil
}

// updates the mod time for a TRACKED file in an index
func update_index_file_mod(index string, filename string, modtime string) error {
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

	records = append(records, []string{filename, modtime})

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
