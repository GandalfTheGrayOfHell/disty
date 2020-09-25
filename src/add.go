package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// `add` updates the Mod time in the index

func add(pwd string, filename string) {
	// check if file exists otherwise panic
	file_info, err := os.Stat(filepath.Join(pwd, filename))
	if err != nil {
		panic("[ERROR] Could not find specified file: " + filename)
	}

	// if current mod time > index mod time then update index
	file_mod := int(file_info.ModTime().Unix())
	s_index_mod, err := get_index_file_modtime(filepath.Join(pwd, ".disty", "index.csv"), filename)
	if err != nil {
		panic("[ERROR] Could not get file modification time")
	}

	index_mod, err := strconv.Atoi(s_index_mod)
	if err != nil {
		panic("[ERROR] Could not parse Mod time")
	}

	if file_mod > index_mod {
		err := update_index_file_mod(filepath.Join(pwd, ".disty", "index.csv"), filename, strconv.Itoa(file_mod))
		if err != nil {
			panic("[ERROR] Could not update mod time for file: " + filepath.Join(pwd, ".disty", "index.csv"))
		}
	}
}

func Add(filenames []string) {
	if len(filenames) < 1 {
		panic("[ERROR] Pass a file to track")
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic("[ERROR] Could not get working directory")
	}

	if filenames[0] == "." {
		// run normal routine
		filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic("[ERROR] `Walk` error for file: " + pwd)
			}

			if !info.IsDir() {
				add(pwd, strings.Replace(path, pwd, "", 1))
			}

			return nil
		})
	} else {
		for _, f := range filenames {
			add(pwd, filepath.FromSlash(f))
		}
	}
}
