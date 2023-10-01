package install

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strings"
)

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

func installAll(commander Commander) {

	var notInstalled []string
	var commands []string

	all := []Dependency{JqDependency(), KarabinerDependency(), AltTabDependency(), RectangleDependency()}

	for _, d := range all {
		if !commander.Exists(d.command) {
			notInstalled = append(notInstalled, d.name)
			commands = append(commands, d.installCommand)
		}
	}

	if len(notInstalled) > 0 {
		installApp := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("The following dependencies will be installed: %s. Do you agree?", strings.Join(notInstalled, ", ")),
		}
		handleInterrupt(survey.AskOne(prompt, &installApp))

		if !installApp {
			fmt.Printf("Qutting...")
			commander.Exit(0)
		}

		for _, c := range commands {
			commander.Run(c)
		}
	}
}
