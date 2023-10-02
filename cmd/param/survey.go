package param

import (
	"github.com/AlecAivazis/survey/v2"
)

var questions = []*survey.Question{
	{
		Name: "appLauncher",
		Prompt: &survey.Select{
			Message: "Your App Launcher (Win/Opt key):",
			Options: []string{Spotlight, Launchpad, Alfred, "None"},
			Help:    "Select you application launcher which will be available under Win/Opt.",
		},
	},
	{
		Name: "terminal",
		Prompt: &survey.Select{
			Message: "Your Terminal (Ctrl+Alt+T/Ctrl+Cmd+T shortcut):",
			Options: []string{Default, ITerm, Warp, "None"},
			Help:    "On Linux systems Ctrl+Alt+T typically starts the default terminal. Let me help with that",
		},
	},
	{
		Name: "keyboardLayout",
		Prompt: &survey.Select{
			Message: "Your external keyboard layout:",
			Options: []string{PC, Mac, "None"},
			Help:    `The layout of your external keyboard to help adjust the setup. If you do not use any, just select "None"`,
		},
	},
	{
		Name: "ides",
		Prompt: &survey.MultiSelect{
			Message: "IDE keymaps to install:",
			Options: IdeKeymapOptions(),
			Help:    "IDEs/tools to apply the PC keymaps for",
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
			Help: "Apps to be blacklisted e.g. they won't be appearing in the Alt-Tab switcher",
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
			Help:     "Additional macOS settings to make your life better",
			PageSize: 15,
		},
	},
}
