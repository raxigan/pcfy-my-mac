package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type Commander interface {
	run(command string)
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
		switch os := runtime.GOOS; os {
		case "darwin":
			BasicCommander{}.run(command)
		case "linux":
			fmt.Println("Running: " + command)
		default:
			fmt.Printf("Cannot run on OS: %s\n", os)
		}
	}
}
