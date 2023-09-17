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
	name              string
	fullName          string
	toolboxScriptName string
	flag              string
}

func IntelliJ() IDE {
	return IDE{
		name:              "IntelliJ",
		fullName:          "IntelliJ IDEA Ultimate",
		toolboxScriptName: "idea",
		flag:              "intellij",
	}
}

func PyCharm() IDE {
	return IDE{
		name:              "PyCharm CE",
		fullName:          "PyCharm Community Edition",
		toolboxScriptName: "pycharm",
		flag:              "pycharm-ce",
	}
}

func GoLand() IDE {
	return IDE{
		name:              "GoLand",
		fullName:          "GoLand",
		toolboxScriptName: "goland",
		flag:              "goland",
	}
}

func Fleet() IDE {
	return IDE{
		name:              "Fleet",
		fullName:          "Fleet",
		toolboxScriptName: "fleet",
	}
}
