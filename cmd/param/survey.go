package param

import (
	"github.com/AlecAivazis/survey/v2"
)

var questions = []*survey.Question{
	{
		Name: "appLauncher",
		Prompt: &survey.Select{
			Message: "Assign Win/Opt key action (open app launcher):",
			Options: []string{Spotlight, Launchpad, Alfred, None},
			Help:    `Select you application launcher which will be available under Win/Opt. Select "None" if you don't use any'`,
		},
	},
	{
		Name: "terminal",
		Prompt: &survey.Select{
			Message: "Assign Ctrl+Alt+T/Ctrl+Cmd+T shortcut action (open terminal):",
			Options: []string{Default, ITerm, Warp, Wave, None},
			Help:    `On Linux systems Ctrl+Alt+T starts the default terminal. Let me take care of that or select "None"`,
		},
	},
	{
		Name: "keyboardLayout",
		Prompt: &survey.Select{
			Message: "Specify the layout of your external keyboard (if any):",
			Options: []string{PC, Mac, None},
			Help:    `The layout of your external keyboard to help adjust the setup. If you do not use any, just select "None"`,
		},
	},
	{
		Name: "systemSettings",
		Prompt: &survey.MultiSelect{
			Message: "Select additional system settings to apply:",
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
			Help: `
Additional macOS settings to make your life better

• Enable Dock auto-hide (2s delay) - partially disable Dock
• Change Dock minimize animation to "scale" - if you don't like animations
• Enable Home & End keys - they have no action assigned by default
• Show hidden files in Finder - always show dot-files
• Show directories on top in Finder - show directories on top
• Show full POSIX paths in Finder - show full path instead of current directory name
`,
			PageSize: 15,
		},
	},
}
