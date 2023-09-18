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
	directory         string
	fullName          string
	toolboxScriptName string
	flag              string
}

func IntelliJ() IDE {
	return IDE{
		directory:         "IntelliJIdea",
		fullName:          "IntelliJ IDEA Ultimate",
		toolboxScriptName: "idea",
		flag:              "intellij",
	}
}

func PyCharm() IDE {
	return IDE{
		directory:         "PyCharmCE",
		fullName:          "PyCharm Community Edition",
		toolboxScriptName: "pycharm",
		flag:              "pycharm-ce",
	}
}

func GoLand() IDE {
	return IDE{
		directory:         "GoLand",
		fullName:          "GoLand",
		toolboxScriptName: "goland",
		flag:              "goland",
	}
}

func Fleet() IDE {
	return IDE{
		directory:         "Fleet",
		fullName:          "Fleet",
		toolboxScriptName: "fleet",
		flag:              "fleet",
	}
}
