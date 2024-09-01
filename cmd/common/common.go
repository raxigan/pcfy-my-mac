package common

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"os"
	"os/exec"
	"strings"
)

const Yellow = "\033[33m"
const Red = "\033[31m"
const Blue = "\033[34m"
const Green = "\033[32m"
const Cyan = "\033[36m"
const Purple = "\033[35m"
const Reset = "\033[0m"

var ExecCommand = exec.Command

func Exists(command string) bool {
	if strings.HasSuffix(command, ".app") {
		cmd := ExecCommand("mdfind", "-name", command)
		output, _ := cmd.Output()
		return strings.TrimSpace(string(output)) != ""
	} else {
		_, err := exec.LookPath(command)
		return err == nil
	}
}

func GetOrDefaultString(launcher string, launcher2 *string) string {
	if launcher2 != nil {
		return *launcher2
	} else {
		return launcher
	}
}

func GetOrDefaultSlice(launcher []string, launcher2 *[]string) []string {
	if launcher2 != nil {
		return *launcher2
	} else {
		return launcher
	}
}

func PrintColored(color, msg string) {
	fmt.Print(Colored(color, msg))
}

func Colored(color, msg string) string {
	return color + msg + Reset
}

func HandleInterrupt(err error) {
	if errors.Is(err, terminal.InterruptErr) {
		fmt.Println("Quitting...")
		os.Exit(1)
	}
}
