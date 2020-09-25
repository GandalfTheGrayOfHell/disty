package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Status() {
	pwd, err := os.Getwd()
	if err != nil {
		panic("[ERROR] Could not get working directory")
	}

	reader, err := ioutil.ReadFile(filepath.Join(pwd, ".disty", "index.csv"))
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

	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		rel_path := strings.Replace(path, pwd, "", 1)

		if err != nil {
			panic("[ERROR] `Walk` error for file: " + rel_path)
		}

		if info.IsDir() && info.Name() == ".disty" {
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
		fmt.Println("[INFO] No changes")
	} else {
		fmt.Println("Modified files:")

		for _, f := range modified {
			fmt.Printf("\t%s\n", f)
		}

		fmt.Println("\nUntracked files:")

		for _, f := range untracked {
			fmt.Printf("\t%s\n", f)
		}

	}
}
