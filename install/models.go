package install

import (
	"errors"
)

type AppLauncher int

const (
	Spotlight AppLauncher = iota
	Launchpad
	Alfred
)

type KeyboardLayout int

const (
	PC KeyboardLayout = iota
	Mac
)

type Terminal int

const (
	Default Terminal = iota
	iTerm
	Warp
)

func (me AppLauncher) String() string {
	return [...]string{"Spotlight", "Launchpad", "Alfred"}[me]
}

func (me KeyboardLayout) String() string {
	return [...]string{"PC", "Mac"}[me]
}

func (me Terminal) String() string {
	return [...]string{"Default", "iTerm", "Warp"}[me]
}

type MySurvey struct {
	Message     string
	Options     []string
	Description func(value string, index int) string
}

type IDE struct {
	parentDir       string // relative to homedir
	dir             string // relative to parentDir
	keymapsDir      string // relative to dir
	fullName        string
	srcKeymapsFile  string
	destKeymapsFile string
	multipleDirs    bool
}

func IntelliJ() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "IntelliJ",
		keymapsDir:      "/keymaps",
		fullName:        "IntelliJ IDEA Ultimate",
		srcKeymapsFile:  "idea.xml",
		destKeymapsFile: "intellij-idea-ultimate.xml",
		multipleDirs:    true,
	}
}

func IntelliJCE() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "IdeaIC",
		keymapsDir:      "/keymaps",
		fullName:        "IntelliJ IDEA CE",
		srcKeymapsFile:  "idea.xml",
		destKeymapsFile: "intellij-idea-community-edition.xml",
		multipleDirs:    true,
	}
}

func PyCharm() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "PyCharmCE",
		keymapsDir:      "/keymaps",
		fullName:        "PyCharm CE",
		srcKeymapsFile:  "idea.xml",
		destKeymapsFile: "pycharm-community-edition.xml",
		multipleDirs:    true,
	}
}

func GoLand() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "GoLand",
		keymapsDir:      "/keymaps",
		fullName:        "GoLand",
		srcKeymapsFile:  "idea.xml",
		destKeymapsFile: "goland.xml",
		multipleDirs:    true,
	}
}

func Fleet() IDE {
	return IDE{
		parentDir:       "",
		dir:             ".fleet",
		keymapsDir:      "/keymap",
		fullName:        "Fleet",
		srcKeymapsFile:  "fleet.json",
		destKeymapsFile: "user.json",
	}
}

var IDEKeymaps = []IDE{IntelliJ(), IntelliJCE(), PyCharm(), GoLand(), Fleet()}
var SystemSettings = []string{
	"Enable Dock auto-hide (2s delay)",
	`Change Dock minimize animation to "scale"`,
	"Enable Home & End keys",
	"Show hidden files in Finder",
	"Show directories on top in Finder",
	"Show full POSIX paths in Finder window title",
}

func IdeKeymapsSurveyOptions() []string {

	var options []string

	for _, e := range IDEKeymaps {
		options = append(options, e.fullName)
	}

	return options
}

func IdeKeymapOptions() []string {

	var options []string

	for _, e := range IDEKeymaps {
		options = append(options, e.fullName)
	}

	return options
}

func IdeKeymapByFullName(fullName string) (IDE, error) {

	for _, e := range IDEKeymaps {
		if e.fullName == fullName {
			return e, nil
		}
	}

	return IDE{}, errors.New("No keymap found: " + fullName)
}
