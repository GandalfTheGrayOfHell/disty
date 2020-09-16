package main

import (
	"flag"
)

func main() {
	username := flag.String("username", "", "Username for Disty")
	password := flag.String("password", "", "Password for Disty")
	flag.Parse()

	command := flag.Arg(0)

	switch command {
	case "init":
		Init()
	case "serve":
		Serve()
	case "config":
		Config(*username, *password)
	}
}
