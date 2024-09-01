package install

import (
	"errors"
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/schollz/progressbar/v3"
	"os"
	"os/exec"
	"strings"
	"time"
)

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
	Exit(code int)
	TryLog(msgType LogMessage, text string)
}

type DefaultCommander struct {
	Verbose  bool
	Progress *progressbar.ProgressBar
}

func NewDefaultCommander(verbose bool) *DefaultCommander {
	return &DefaultCommander{
		Verbose:  verbose,
		Progress: progressbar.NewOptions(100, progressbar.OptionSetWidth(60)),
	}
}

func (c *DefaultCommander) Run(command string) {

	if command == "clear" {
		if !c.Verbose {
			clearConsole()
		}
		return
	}

	c.TryLog(CmdMsg, command)

	var out []byte = nil
	var err error = nil

	if strings.HasPrefix(command, "jq") {
		out, err = exec.Command("/bin/bash", "-c", command).CombinedOutput()
	} else {
		out, err = common.ExecCommand("/bin/bash", "-c", command).CombinedOutput()
	}

	if !strings.HasPrefix(command, "killall") && err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			c.TryLog(StdErrMsg, string(out))
			os.Exit(1)
		}
	}

	c.TryLog(StdOutMsg, string(out))
}

func clearConsole() {
	cmd := common.ExecCommand("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (c *DefaultCommander) Exit(code int) {
	os.Exit(code)
}

var (
	TaskMsg   = LogMessage{"TASK", common.Blue}
	CmdMsg    = LogMessage{"CMD", common.Green}
	WarnMsg   = LogMessage{"WARN", common.Yellow}
	StdOutMsg = LogMessage{"STDOUT", common.Purple}
	FileMsg   = LogMessage{"FILE", common.Cyan}

	ErrMsg    = LogMessage{"ERROR", common.Red}
	StdErrMsg = LogMessage{"STDERR", common.Red}
)

type LogMessage struct {
	msgType string
	color   string
}

func (c *DefaultCommander) TryLog(logMsg LogMessage, output string) {

	output = strings.ReplaceAll(output, os.Getenv("HOME"), "~")

	if logMsg.msgType == ErrMsg.msgType || logMsg.msgType == StdErrMsg.msgType {
		log(logMsg, output)
	} else if !c.Verbose && c.Progress != nil {
		c.Progress.Add(4)
	} else if len(output) != 0 {
		log(logMsg, output)
	}
}

func log(msg LogMessage, output string) {

	isHelper := os.Getenv("GO_WANT_HELPER_PROCESS")

	if isHelper == "" {
		fmt.Println(fmt.Sprintf("[%s] %s", common.Colored(msg.color, msg.msgType), output))
	} else {
		fmt.Println(output)
	}
}
