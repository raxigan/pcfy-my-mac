package install_test

import (
	"fmt"
	"github.com/raxigan/pcfy-my-mac/install"
	"os"
	"strings"
	"time"
)

type MockCommander struct {
	CommandsLog []string
}

func (c *MockCommander) Run(command string) {

	cmd := strings.Fields(command)[0]

	fmt.Println("Running: " + command)

	switch cmd {
	case "jq":
		install.DefaultCommander{}.Run(command)
	case "killall", "open", "clear", "defaults":
		c.CommandsLog = append(c.CommandsLog, command)
	case "plutil":
		pwd, _ := os.Getwd()
		c.CommandsLog = append(c.CommandsLog, strings.ReplaceAll(command, pwd, ""))
	default:
		fmt.Println("Cannot execute command: " + command)
		os.Exit(1)
	}
}

func (c *MockCommander) Exists(command string) bool {
	return true
}

type FakeTimeProvider struct {
}

func (tp FakeTimeProvider) Now() time.Time {
	parse, _ := time.Parse("2006-01-02 15:04:05", "2023-09-27 12:30:00")
	return parse
}
