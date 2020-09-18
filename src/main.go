package main

import (
	"flag"
)

func main() {
	// TODO: take username and password from stdin
	username := flag.String("username", "", "Username for Disty")
	password := flag.String("password", "", "Password for Disty")
	remote := flag.String("url", "", "URL for Remote server")
	flag.Parse()

	command := flag.Arg(0)

	switch command {
	case "init":
		Init()
	case "serve":
		Serve()
	case "config":
		Config(*username, *password)
	case "remote":
		Remote(*remote)
	}
}
