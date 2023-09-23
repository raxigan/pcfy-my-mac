package main

import (
	"errors"
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

	out, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Printf("Command failed with error: %s\n", exitErr)
			fmt.Printf("Stderr: %s", out)
		}
		return
	}

	fmt.Print(string(out))
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
