package install_test

import (
	"fmt"
	"github.com/raxigan/pcfy-my-mac/install"
	"os"
	"strings"
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
