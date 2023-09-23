package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const YELLOW = "\x1b[33m%s\x1b[0m"

type Commander interface {
	run(command string)
	exists(command string) bool
}

type BasicCommander struct {
}

func (c BasicCommander) run(command string) {
	fmt.Println("Running: " + command)

	out, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()

	if strings.Fields(command)[0] != "killall" && err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Printf("Command failed with error: %s\n", exitErr)
			fmt.Printf("Stderr: %s", out)
			os.Exit(1)
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
	case "jq", "plutil", "defaults":
		BasicCommander{}.run(command)
	case "killall", "open", "clear":
		fmt.Println("Running: " + command)
	default:
		fmt.Println("Cannot execute command: " + command)
		os.Exit(1)
	}
}

func (c MockCommander) exists(command string) bool {
	return true
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func printColored(color, msg string) {
	fmt.Println(fmt.Sprintf(color, msg))
}
