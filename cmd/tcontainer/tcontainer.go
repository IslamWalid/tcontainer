package main

import (
	"fmt"
	"os"

	"github.com/IslamWalid/tcontainer/internal/container"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "too few arguments")
		os.Exit(1)
	}

	err := container.Initialize()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		err := container.Run(os.Args[2], os.Args[3:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	case "child":
		err := container.Child(os.Args[2], os.Args[3:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}
