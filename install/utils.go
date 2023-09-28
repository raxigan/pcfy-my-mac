package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const YELLOW = "\x1b[33m%s\x1b[0m"

type TimeProvider interface {
	Now() time.Time
}

type DefaultTimeProvider struct {
}

func (tp DefaultTimeProvider) Now() time.Time {
	return time.Now()
}

type Commander interface {
	Run(command string)
	Exists(command string) bool
}

type DefaultCommander struct {
}

func (c DefaultCommander) Run(command string) {
	fmt.Println("Running: " + command)

	out, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()

	if strings.Fields(command)[0] != "killall" && err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Printf("Command failed with error: %s\n", exitErr)
			fmt.Printf("Stderr: %s", out)
			os.Exit(1)
		}
	}

	fmt.Print(string(out))
}

func (c DefaultCommander) Exists(command string) bool {
	if strings.HasSuffix(command, ".app") {
		return fileExists(command)
	} else {
		_, err := exec.LookPath(command)
		return err == nil
	}
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
