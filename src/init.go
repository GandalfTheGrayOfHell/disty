package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Init() {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	// root dir
	disty := pwd + "/.disty"
	if os.MkdirAll(disty, 0777) != nil {
		return
	}

	index, err := os.Create(disty + "/index.csv")
	defer index.Close()
	if err != nil {
		return
	}

	csvwriter := csv.NewWriter(index)
	defer csvwriter.Flush()

	// make a copy of all the files
	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("%q: %v", path, err)
			return err
		}

		// ignore the disty directory completely
		if info.IsDir() && info.Name() == ".disty" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			mod_time := strconv.FormatInt(info.ModTime().Unix(), 10) // modification time
			rel_path := strings.Replace(path, pwd, "", 1)            // path with pre-root removed
			csvwriter.Write([]string{rel_path, mod_time})
		}

		return nil
	})
}
