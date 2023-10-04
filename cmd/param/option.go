package param

import (
	"errors"
	"strings"
)

const (
	Spotlight = "Spotlight"
	Launchpad = "Launchpad"
	Alfred    = "Alfred"

	Default = "Default"
	ITerm   = "iTerm"
	Warp    = "Warp"

	PC   = "PC"
	Mac  = "Mac"
	None = "None"

	Spotify           = "Spotify"
	Finder            = "Finder"
	SystemPreferences = "System Preferences"
	AltTab            = "AltTab"
)

var AppToBundleMapping = map[string]string{
	strings.ToLower(Spotify):           "com.spotify.client",
	strings.ToLower(Finder):            "com.apple.finder",
	strings.ToLower(SystemPreferences): "com.apple.systempreferences",
	strings.ToLower(ITerm):             "com.googlecode.iterm2",
	strings.ToLower(AltTab):            "com.lwouis.alt-tab-macos",
}

type IDE struct {
	ParentDir       string // relative to homedir
	Dir             string // relative to parentDir
	KeymapsDir      string // relative to dir
	FullName        string
	SrcKeymapsFile  string
	DestKeymapsFile string
	MultipleDirs    bool
}

func IntelliJ() IDE {
	return IDE{
		ParentDir:       "Library/Application Support/JetBrains",
		Dir:             "IntelliJ",
		KeymapsDir:      "keymaps",
		FullName:        "IntelliJ IDEA Ultimate",
		SrcKeymapsFile:  "idea.xml",
		DestKeymapsFile: "intellij-idea-ultimate.xml",
		MultipleDirs:    true,
	}
}

func IntelliJCE() IDE {
	return IDE{
		ParentDir:       "Library/Application Support/JetBrains",
		Dir:             "IdeaIC",
		KeymapsDir:      "keymaps",
		FullName:        "IntelliJ IDEA CE",
		SrcKeymapsFile:  "idea.xml",
		DestKeymapsFile: "intellij-idea-community-edition.xml",
		MultipleDirs:    true,
	}
}

func PyCharm() IDE {
	return IDE{
		ParentDir:       "Library/Application Support/JetBrains",
		Dir:             "PyCharmCE",
		KeymapsDir:      "keymaps",
		FullName:        "PyCharm CE",
		SrcKeymapsFile:  "idea.xml",
		DestKeymapsFile: "pycharm-community-edition.xml",
		MultipleDirs:    true,
	}
}

func GoLand() IDE {
	return IDE{
		ParentDir:       "Library/Application Support/JetBrains",
		Dir:             "GoLand",
		KeymapsDir:      "keymaps",
		FullName:        "GoLand",
		SrcKeymapsFile:  "idea.xml",
		DestKeymapsFile: "goland.xml",
		MultipleDirs:    true,
	}
}

func Fleet() IDE {
	return IDE{
		ParentDir:       "",
		Dir:             ".fleet",
		KeymapsDir:      "keymap",
		FullName:        "Fleet",
		SrcKeymapsFile:  "fleet.json",
		DestKeymapsFile: "user.json",
	}
}

var IDEKeymaps = []IDE{IntelliJ(), IntelliJCE(), PyCharm(), GoLand(), Fleet()}
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
		if strings.ToLower(e.FullName) == strings.ToLower(fullName) {
			return e, nil
		}
	}

	return IDE{}, errors.New("No keymap found: " + fullName)
}
