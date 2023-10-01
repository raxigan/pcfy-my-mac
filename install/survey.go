package install

import "github.com/AlecAivazis/survey/v2"

var questions = []*survey.Question{
	{
		Name: "appLauncher",
		Prompt: &survey.Select{
			Message: "Your App Launcher (Win/Opt key):",
			Options: []string{Spotlight.String(), Launchpad.String(), Alfred.String(), "None"},
		},
	},
	{
		Name: "terminal",
		Prompt: &survey.Select{
			Message: "Your Terminal (Ctrl+Alt+T/Ctrl+Cmd+T shortcut):",
			Options: []string{Default.String(), iTerm.String(), Warp.String(), "None"},
		},
	},
	{
		Name: "keyboardLayout",
		Prompt: &survey.Select{
			Message: "Your external keyboard layout:",
			Options: []string{PC.String(), Mac.String(), "None"},
		},
	},
	{
		Name: "ides",
		Prompt: &survey.MultiSelect{
			Message: "IDE keymaps to install:",
			Options: IdeKeymapOptions(),
		},
	},
	{
		Name: "blacklist",
		Prompt: &survey.MultiSelect{
			Message: "Apps to blacklist:",
			Options: []string{
				"Spotify",
				"Finder",
				"System Preferences",
			},
			Help: "help",
		},
	},
	{
		Name: "systemSettings",
		Prompt: &survey.MultiSelect{
			Message: "System settings:",
			Options: []string{
				"Enable Dock auto-hide (2s delay)",
				`Change Dock minimize animation to "scale"`,
				"Enable Home & End keys",
				"Show hidden files in Finder",
				"Show directories on top in Finder",
				"Show full POSIX paths in Finder",
			},
			Description: func(value string, index int) string {
				if index < 2 {
					return "Recommended"
				}
				return ""
			},
			Help:     "help",
			PageSize: 15,
		},
	},
}
