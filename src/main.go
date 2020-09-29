package main

import (
	"flag"
	"os"
)

func main() {
	// TODO: take username and password from stdin
	// initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	// pushCmd := flag.NewFlagSet("push", flag.ExitOnError)
	// pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	// statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	// addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	cloneCmd := flag.NewFlagSet("clone", flag.ExitOnError)
	cloneUrl := cloneCmd.String("url", "", "URL for Remote server")
	cloneDir := cloneCmd.String("dir", "./", "Directory path to clone repo")
	cloneProject := cloneCmd.String("project", "", "Project name")

	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	configUser := configCmd.String("username", "", "Username for Disty")
	configPass := configCmd.String("password", "", "Password for Disty")

	remoteCmd := flag.NewFlagSet("remote", flag.ExitOnError)
	remoteServer := remoteCmd.String("url", "", "URL for Remote server")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	servePort := serveCmd.Int("port", 3000, "Port to start Disty server")
	serveDir := serveCmd.String("dir", "./", "Directory path for server to store projects on")

	switch os.Args[1] {
	case "init":
		Init()
	case "serve":
		serveCmd.Parse(os.Args[2:])
		Serve(*servePort, *serveDir)
	case "config":
		Config(*configUser, *configPass)
	case "remote":
		remoteCmd.Parse(os.Args[2:])
		Remote(*remoteServer)
	case "push":
		Push()
	case "clone":
		cloneCmd.Parse(os.Args[2:])
		Clone(*cloneUrl, *cloneDir, *cloneProject)
	case "status":
		Status()
	case "add":
		Add(os.Args[2:])
	}
}
