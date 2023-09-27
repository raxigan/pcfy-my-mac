package main

import (
	"log"
	"os"
)
import "github.com/raxigan/pcfy-my-mac/install"

func main() {
	homeDir, _ := os.UserHomeDir()
	_, err := install.RunInstaller(homeDir, install.DefaultCommander{}, nil)

	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	log.Fatalf("Error: %s", err)
}
