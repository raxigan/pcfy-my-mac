package task

type Dependency struct {
	name           string
	command        string
	installCommand string
}

func KarabinerDependency() Dependency {
	return Dependency{
		name:           "Karabiner-Elements",
		command:        "Karabiner-Elements.app",
		installCommand: "brew install --cask karabiner-elements",
	}
}

func JqDependency() Dependency {
	return Dependency{
		name:           "jq",
		command:        "jq",
		installCommand: "brew install jq",
	}
}

func AltTabDependency() Dependency {
	return Dependency{
		name:           "AltTab",
		command:        "AltTab.app",
		installCommand: "brew install --cask alt-tab",
	}
}

func RectangleDependency() Dependency {
	return Dependency{
		name:           "Rectangle",
		command:        "Rectangle.app",
		installCommand: "brew install --cask rectangle",
	}
}
