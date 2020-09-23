package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

// Module should just allow addition of username and password of an individual user

func Config(username, password string) {
	if username == "" || password == "" {
		panic("[ERROR] Username and password is required")
	}

	if os.MkdirAll(path.Join(os.TempDir(), "disty"), 0777) != nil {
		panic("[ERROR] Could not create dir: " + os.TempDir() + "\\disty")
	}

	f, err := os.Create(path.Join(os.TempDir(), "disty", "auth"))
	defer f.Close()
	if err != nil {
		panic("[ERROR] Could not create file: " + os.TempDir() + "\\disty\\auth")
	}

	md5_user := md5.Sum([]byte(username))
	md5_pass := md5.Sum([]byte(password))

	hex_user := hex.EncodeToString(md5_user[:])
	hex_pass := hex.EncodeToString(md5_pass[:])

	f.Write([]byte(hex_user))
	f.Write([]byte("\n"))
	f.Write([]byte(hex_pass))

	fmt.Println("[SUCCESS] Configuration saved at", os.TempDir()+"\\disty")
}
