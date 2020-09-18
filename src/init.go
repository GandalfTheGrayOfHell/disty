package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO: add project.json
func Init() {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	// root dir
	disty := pwd + "\\.disty"
	if os.MkdirAll(disty, 0777) != nil {
		return
	}

	project_file, err := os.Create(disty + "\\project.json")
	defer project_file.Close()
	if err != nil {
		return
	}

	project := Project{}
	project.Default()

	project_json, err := json.MarshalIndent(project, "", "\t")
	if err != nil {
		return
	}

	project_file.Write(project_json)

	index_file, err := os.Create(disty + "\\index.csv")
	defer index_file.Close()
	if err != nil {
		return
	}

	csvwriter := csv.NewWriter(index_file)
	defer csvwriter.Flush()

	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("%q: %v", path, err)
			return err
		}

		if !info.IsDir() {
			mod_time := strconv.FormatInt(info.ModTime().Unix(), 10) // modification time
			rel_path := strings.Replace(path, pwd, "", 1)            // path with pre-root removed
			csvwriter.Write([]string{rel_path, mod_time})
		}

		return nil
	})

	fmt.Println("[SUCCESS] Initialized successfully")
}
