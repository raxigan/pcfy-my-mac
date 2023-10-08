package install

import (
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"time"
)

type Installation struct {
	Commander
	param.Params
	HomeDir
	ProfileName      string
	InstallationTime time.Time
}
