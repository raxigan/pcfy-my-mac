package param

import (
	"errors"
)

const (
	Spotlight = "Spotlight"
	Launchpad = "Launchpad"
	Alfred    = "Alfred"

	Default = "Default"
	ITerm   = "iTerm"
	Warp    = "Warp"
	Wave    = "Wave"

	PC   = "PC"
	Mac  = "Mac"
	None = "None"
)

type IDE struct {
	KeymapsDir      string // relative to homedir
	FullName        string
	SrcKeymapsFile  string
	DestKeymapsFile string
}

func IntelliJ() IDE {
	return IDE{
		KeymapsDir:     "Library/Application Support/JetBrains/IntelliJ{version}/keymaps",
		FullName:       "IntelliJ IDEA Ultimate",
		SrcKeymapsFile: "idea.xml",
	}
}

func IntelliJCE() IDE {
	return IDE{
		KeymapsDir:     "Library/Application Support/JetBrains/IdeaIC{version}/keymaps",
		FullName:       "IntelliJ IDEA Community Edition",
		SrcKeymapsFile: "idea.xml",
	}
}

func PyCharmCE() IDE {
	return IDE{
		KeymapsDir:     "Library/Application Support/JetBrains/PyCharmCE{version}/keymaps",
		FullName:       "PyCharm Community Edition",
		SrcKeymapsFile: "idea.xml",
	}
}

func PyCharm() IDE {
	return IDE{
		KeymapsDir:     "Library/Application Support/JetBrains/PyCharm{version}/keymaps",
		FullName:       "PyCharm Professional Edition",
		SrcKeymapsFile: "idea.xml",
	}
}

func GoLand() IDE {
	return IDE{
		KeymapsDir:     "Library/Application Support/JetBrains/GoLand{version}/keymaps",
		FullName:       "GoLand",
		SrcKeymapsFile: "idea.xml",
	}
}

func AndroidStudio() IDE {
	return IDE{
		KeymapsDir:     "Library/Application Support/Google/AndroidStudio{version}/keymaps",
		FullName:       "Android Studio",
		SrcKeymapsFile: "idea.xml",
	}
}

func Fleet() IDE {
	return IDE{
		KeymapsDir:      ".fleet/keymap",
		FullName:        "Fleet",
		SrcKeymapsFile:  "fleet.json",
		DestKeymapsFile: "user.json",
	}
}

var IDEKeymaps = []IDE{IntelliJ(), IntelliJCE(), PyCharm(), PyCharmCE(), GoLand(), AndroidStudio(), Fleet()}
var SystemSettings = []string{
	"Enable Dock auto-hide (2s delay)",
	`Change Dock minimize animation to "scale"`,
	"Enable Home and End keys",
	"Show hidden files in Finder",
	"Show directories on top in Finder",
	"Show full POSIX paths in Finder window title",
}

func IdeKeymapOptions() []string {

	var options []string

	for _, e := range IDEKeymaps {
		options = append(options, e.FullName)
	}

	return options
}

func IdeKeymapByFullName(fullName string) (IDE, error) {

	for _, e := range IDEKeymaps {
		if ToSimpleParamName(e.FullName) == ToSimpleParamName(fullName) {
			return e, nil
		}
	}

	return IDE{}, errors.New("No keymap found: " + fullName)
}
