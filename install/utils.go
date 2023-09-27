package install

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

type DefaultCommander struct {
}

func (c DefaultCommander) run(command string) {
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

func (c DefaultCommander) exists(command string) bool {
	if strings.HasSuffix(command, ".app") {
		return fileExists(command)
	} else {
		_, err := exec.LookPath(command)
		return err == nil
	}
}

type MockCommander struct {
	CommandsLog []string
}

func (c *MockCommander) run(command string) {

	cmd := strings.Fields(command)[0]

	switch cmd {
	case "jq", "plutil":
		DefaultCommander{}.run(command)
	case "killall", "open", "clear", "defaults":
		fmt.Println("Running: " + command)
		c.CommandsLog = append(c.CommandsLog, command)
	default:
		fmt.Println("Cannot execute command: " + command)
		os.Exit(1)
	}
}

func (c *MockCommander) exists(command string) bool {
	return true
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func printColored(color, msg string) {
	fmt.Println(fmt.Sprintf(color, msg))
}

func replaceWordInFile(path, oldWord, newWord string) error {

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	modifiedContent := strings.ReplaceAll(string(content), oldWord, newWord)

	err = os.WriteFile(path, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
