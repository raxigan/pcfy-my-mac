package main

import "os"
import "github.com/raxigan/pcfy-my-mac/install"

func main() {
	homeDir, _ := os.UserHomeDir()
	install.RunInstaller(homeDir, install.DefaultCommander{}, nil)
}
