// go build -o disty.exe ./src && pushd test && ..\disty.exe init && popd
package main

import (
	"flag"
)

func main() {
	flag.Parse()

	command := flag.Arg(0) // serve, push, pull, init

	switch command {
	case "init":
		{
			i := Init{}
			i.New()
		}
	}
}
