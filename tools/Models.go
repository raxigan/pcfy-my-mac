package main

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
	description string
	options     []string
}

type IDE struct {
	parentDir         string // relative to homedir
	dir               string // relative to parentDir
	keymapsDir        string // relative to dir
	fullName          string
	toolboxScriptName string
	flag              string
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
		flag:            "intellij",
		srcKeymapsFile:  "intellij-idea-ultimate.xml",
		destKeymapsFile: "intellij-idea-ultimate.xml",
		requiresPlugin:  true,
	}
}

func PyCharm() IDE {
	return IDE{
		parentDir:       "/Library/Application Support/JetBrains/",
		dir:             "PyCharmCE",
		keymapsDir:      "/keymaps",
		fullName:        "PyCharm Community Edition",
		flag:            "pycharm-ce",
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
		flag:            "goland",
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
		flag:            "fleet",
		srcKeymapsFile:  "fleet.json",
		destKeymapsFile: "user.json",
	}
}
