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

func Init(projectname string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic("[ERROR] Could not get working directory")
	}

	// root dir
	disty := filepath.Join(pwd, ".disty")
	if os.MkdirAll(disty, 0777) != nil {
		panic("[ERROR] Could not create dir: " + pwd + "\\.disty")
	}

	project_file, err := os.Create(filepath.Join(disty, "project.json"))
	defer project_file.Close()
	if err != nil {
		panic("[ERROR] Could not create file: " + disty + "\\project.json")
	}

	project := Project{Name: projectname, Remote: ""}

	project_json, err := json.MarshalIndent(project, "", "\t")
	if err != nil {
		panic("[ERROR] Could not JSON Marshal")
	}

	project_file.Write(project_json)

	index_file, err := os.Create(filepath.Join(disty, "index.csv"))
	defer index_file.Close()
	if err != nil {
		panic("[ERROR] Could not create file: " + disty + "\\index.csv")
	}

	csvwriter := csv.NewWriter(index_file)
	defer csvwriter.Flush()

	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic("[ERROR] `Walk` error for file: " + path)
		}

		// ignoring all git files because git handles it better than disty
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			mod_time := strconv.FormatInt(info.ModTime().Unix(), 10) // modification time
			rel_path := strings.Replace(path, pwd, "", 1)            // path with pre-root removed
			csvwriter.Write([]string{rel_path, mod_time, "0"})
		}

		return nil
	})

	fmt.Println("[SUCCESS] Initialized successfully")
}
