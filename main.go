package main

import (
	"fmt"
	"github.com/raxigan/pcfy-my-mac/install"
	"os"
)

func main() {

	handleError(
		install.RunInstaller(install.DefaultHomeDir(), install.DefaultCommander{}, install.DefaultTimeProvider{}, nil),
	)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}
