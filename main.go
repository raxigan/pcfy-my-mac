package main

import "os"
import "github.com/raxigan/macos-pc-mode/install"

func main() {
	homeDir, _ := os.UserHomeDir()
	install.RunInstaller(homeDir, install.DefaultCommander{})
}
