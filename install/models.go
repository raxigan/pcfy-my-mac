package install

import "errors"

type AppLauncher int

const (
	Spotlight AppLauncher = iota
	Launchpad
	Alfred
)

type KeyboardType int

const (
	PC KeyboardType = iota
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

func (me KeyboardType) String() string {
	return [...]string{"PC", "Mac"}[me]
}

func (me Terminal) String() string {
	return [...]string{"Default", "iTerm", "Warp"}[me]
}

type MySurvey struct {
	message     string
	options     []string
	description func(value string, index int) string
}

type IDE struct {
	parentDir         string // relative to homedir
	dir               string // relative to parentDir
	keymapsDir        string // relative to dir
	fullName          string
	toolboxScriptName string
	srcKeymapsFile    string
	destKeymapsFile   string
	requiresPlugin    bool
}

func IntelliJ() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "IntelliJ",
		keymapsDir:      "/keymaps",
		fullName:        "IntelliJ IDEA Ultimate",
		srcKeymapsFile:  "intellij-idea-ultimate.xml",
		destKeymapsFile: "intellij-idea-ultimate.xml",
		requiresPlugin:  true,
	}
}

func IntelliJCE() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "IdeaIC",
		keymapsDir:      "/keymaps",
		fullName:        "IntelliJ IDEA CE",
		srcKeymapsFile:  "intellij-idea-community-edition.xml",
		destKeymapsFile: "intellij-idea-community-edition.xml",
		requiresPlugin:  true,
	}
}

func PyCharm() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "PyCharmCE",
		keymapsDir:      "/keymaps",
		fullName:        "PyCharm CE",
		srcKeymapsFile:  "pycharm-community-edition.xml",
		destKeymapsFile: "pycharm-community-edition.xml",
		requiresPlugin:  true,
	}
}

func GoLand() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "GoLand",
		keymapsDir:      "/keymaps",
		fullName:        "GoLand",
		srcKeymapsFile:  "goland.xml",
		destKeymapsFile: "goland.xml",
		requiresPlugin:  true,
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
var AdditionalOptions = []string{
	"Enable Dock auto-hide (2s delay)",
	`Change Dock minimize animation to "scale"`,
	"Enable Home & End keys",
	"Show hidden files in Finder",
	"Show directories on top in Finder",
	"Show full POSIX paths in Finder",
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
