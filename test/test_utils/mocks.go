package test_utils

import (
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"os"
	"strings"
	"time"
)

type MockCommander struct {
	DefaultCommander *install.DefaultCommander
	CommandsLog      []string
}

func (c *MockCommander) TryPrint(prefix, text string) {
	c.DefaultCommander.TryPrint(prefix, text)
}

func NewMockCommander() *MockCommander {
	return &MockCommander{
		DefaultCommander: &install.DefaultCommander{Verbose: false},
	}
}

func (c *MockCommander) Run(command string) {

	cmd := strings.Fields(command)[0]

	switch cmd {
	case "jq":
		c.DefaultCommander.Run(command)
	case "killall", "open", "clear", "defaults", "brew":
		fmt.Println(fmt.Sprintf("[%s] %s", common.Colored(common.Green, "RUN"), command))
		c.CommandsLog = append(c.CommandsLog, command)
	case "plutil":
		fmt.Println(fmt.Sprintf("[%s] %s", common.Colored(common.Green, "RUN"), command))
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

func (c *MockCommander) Exit(code int) {}

type FakeTimeProvider struct {
}

func (tp FakeTimeProvider) Now() time.Time {
	parse, _ := time.Parse("2006-01-02 15:04:05", "2023-09-27 12:30:00")
	return parse
}
