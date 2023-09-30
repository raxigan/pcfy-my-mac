package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	Exit(code int)
}

type DefaultCommander struct {
	Verbose bool
}

func (c DefaultCommander) Run(command string) {
	c.tryPrint("Running: " + command)

	out, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()

	if strings.Fields(command)[0] != "killall" && err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Printf("Command failed with error: %s\n", exitErr)
			fmt.Printf("Stderr: %s", out)
			os.Exit(1)
		}
	}

	c.tryPrint(string(out))
}

func (c DefaultCommander) Exists(command string) bool {
	if strings.HasSuffix(command, ".app") {
		return fileExists(filepath.Join("/Applications", command))
	} else {
		_, err := exec.LookPath(command)
		return err == nil
	}
}

func (c DefaultCommander) Exit(code int) {
	os.Exit(code)
}

func (c DefaultCommander) tryPrint(output string) {
	if c.Verbose {
		fmt.Print(output)
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

func TextFromFile(paramsFile string) (string, error) {
	d, e := os.ReadFile(paramsFile)

	if e != nil {
		return "", e
	}

	return string(d), nil
}
