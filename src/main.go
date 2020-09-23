package main

import (
	"flag"
)

func main() {
	// TODO: take username and password from stdin
	username := flag.String("username", "", "Username for Disty")
	password := flag.String("password", "", "Password for Disty")
	remote := flag.String("url", "", "URL for Remote server")
	port := flag.Int("port", 3000, "Port to start Disty server")
	dir := flag.String("dir", "./", "Directory path for server to store projects on")
	flag.Parse()

	command := flag.Arg(0)

	switch command {
	case "init":
		Init()
	case "serve":
		Serve(*port, *dir)
	case "config":
		Config(*username, *password)
	case "remote":
		Remote(*remote)
	case "push":
		Push()
	case "clone":
		Clone()
	case "status":
		Status()
	}
}
