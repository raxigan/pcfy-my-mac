package install

import (
	"errors"
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/schollz/progressbar/v3"
	"os"
	"os/exec"
	"path/filepath"
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
	Exists(command string) bool
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

	c.TryLog(CmdMsg, command)

	out, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()

	if strings.Fields(command)[0] != "killall" && err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			c.TryLog(StdErrMsg, string(out))
			os.Exit(1)
		}
	}

	c.TryLog(StdOutMsg, string(out))
}

func (c *DefaultCommander) Exists(command string) bool {
	if strings.HasSuffix(command, ".app") {
		return common.FileExists(filepath.Join("/Applications", command))
	} else {
		_, err := exec.LookPath(command)
		return err == nil
	}
}

func (c *DefaultCommander) Exit(code int) {
	os.Exit(code)
}

var (
	TaskMsg   = LogMessage{"TASK", common.Blue}
	CmdMsg    = LogMessage{"CMD", common.Green}
	WarnMsg   = LogMessage{"WARN", common.Yellow}
	StdOutMsg = LogMessage{"STDOUT", common.Purple}

	ErrMsg    = LogMessage{"ERROR", common.Red}
	StdErrMsg = LogMessage{"STDERR", common.Red}
)

type LogMessage struct {
	msgType string
	color   string
}

func (c *DefaultCommander) TryLog(logMsg LogMessage, output string) {

	if logMsg.msgType == "ERROR" || logMsg.msgType == "STDERR" {
		log(logMsg, output)
	} else if !c.Verbose && c.Progress != nil {
		c.Progress.Add(4)
	} else if len(output) != 0 {
		log(logMsg, output)
	}
}

func log(msg LogMessage, output string) {
	fmt.Println(fmt.Sprintf("[%s] %s", common.Colored(msg.color, msg.msgType), output))
}
