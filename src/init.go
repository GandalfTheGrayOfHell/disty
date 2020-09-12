package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Init struct{}

func (i *Init) New() {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	// root dir
	disty := pwd + "/.disty"
	if os.MkdirAll(disty, 0777) != nil {
		return
	}

	// make a copy of all the files
	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("%q: %v", path, err)
			return err
		}

		if info.IsDir() && info.Name() == ".disty" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			fmt.Println(strings.Replace(path, pwd, "", 1))
		}

		return nil
	})
}
