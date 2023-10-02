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
	TryPrint(prefix, text string)
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

	c.TryPrint(common.Colored(common.Green, "RUN"), command)

	out, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()

	if strings.Fields(command)[0] != "killall" && err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			c.TryPrint(common.Colored(common.Red, "STDERR"), string(out))
			os.Exit(1)
		}
	}

	c.TryPrint(common.Colored(common.Yellow, "STDOUT"), string(out))
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

func (c *DefaultCommander) TryPrint(prefix, output string) {
	if !c.Verbose && c.Progress != nil {
		c.Progress.Add(4)
	} else if len(output) != 0 {
		fmt.Println(fmt.Sprintf("[%s] %s", prefix, output))
	}
}
