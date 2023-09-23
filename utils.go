package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Commander interface {
	run(command string)
	exists(command string) bool
}

type BasicCommander struct {
}

func (c BasicCommander) run(command string) {
	fmt.Println("Running: " + command)

	output, err := exec.Command("/bin/bash", "-c", command).Output()

	if err != nil {
		fmt.Println("Error executing command: "+command+"\n", err)
	}

	fmt.Print(string(output))
}

func (c BasicCommander) exists(command string) bool {
	if strings.HasSuffix(command, ".app") {
		return fileExists(command)
	} else {
		_, err := exec.LookPath(command)
		return err == nil
	}
}

type MockCommander struct {
}

func (c MockCommander) run(command string) {

	cmd := strings.Fields(command)[0]

	switch cmd {
	case "jq", "clean":
		BasicCommander{}.run(command)
	case "killall", "open":
		fmt.Println("Running: " + command)
	case "plutil", "defaults":
		switch system := runtime.GOOS; system {
		case "darwin":
			BasicCommander{}.run(command)
		case "linux":
			fmt.Println("Running: " + command)
		default:
			fmt.Printf("Cannot run on OS: %s\n", system)
		}
	}
}

func (c MockCommander) exists(command string) bool {
	switch system := runtime.GOOS; system {
	case "darwin":
		return BasicCommander{}.exists(command)
	default:
		return true
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
