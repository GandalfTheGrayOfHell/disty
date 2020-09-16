package main

import (
	"fmt"
	"log"
	"os"
)

func Config(username, password string) {
	if username == "" || password == "" {
		fmt.Println(username, password)
		log.Fatalf("[Error] Username and password is required")
	}

	fmt.Println(os.TempDir())
}
