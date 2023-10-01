package install

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Task struct {
	name string
	exec func(i Installation) error
}

func DownloadDependencies() Task {
	return Task{
		name: "Install dependencies",
		exec: func(i Installation) error { installAll(i.Commander); return nil },
	}
}

func CloseKarabiner() Task {
	return Task{
		name: "Close Karabiner",
		exec: func(i Installation) error { i.Run("killall Karabiner-Elements"); return nil },
	}
}

func BackupKarabinerConfig() Task {
	return Task{
		name: "Do karabiner config backup",
		exec: func(i Installation) error {
			original := i.KarabinerConfigFile()
			backupDest := i.KarabinerConfigBackupFile(i.installationTime)
			CopyFile(original, backupDest)
			return nil
		},
	}
}

func DeleteExistingKarabinerProfile() Task {
	return Task{
		name: "Delete existing Karabiner profile",
		exec: func(i Installation) error {
			deleteProfileJqCmd := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >tmp && mv tmp %s", i.profileName, i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(deleteProfileJqCmd)
			return nil
		},
	}
}

func CreateKarabinerProfile() Task {
	return Task{
		name: "Delete existing Karabiner profile",
		exec: func(i Installation) error {
			copyFileFromEmbedFS("karabiner/karabiner-profile.json", "tmp")
			addProfileJqCmd := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile tmp --indent 4 >INPUT.tmp && mv INPUT.tmp %s && rm tmp", i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(addProfileJqCmd)
			return nil
		},
	}
}

func NameKarabinerProfile() Task {
	return Task{
		name: "Delete existing Karabiner profile",
		exec: func(i Installation) error {
			copyFileFromEmbedFS("karabiner/karabiner-profile.json", "tmp")
			addProfileJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name == \"_PROFILE_NAME_\" then .name = \"%s\" else . end)' %s > tmp && mv tmp %s", i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(addProfileJqCmd)
			return nil
		},
	}
}

func UnselectOtherKarabinerProfiles() Task {
	return Task{
		name: "Unselect other Karabiner profiles",
		exec: func(i Installation) error {
			unselectJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name != \"%s\" then .selected = false else . end)' %s > tmp && mv tmp %s", i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(unselectJqCmd)
			return nil
		},
	}
}

func ApplyMainKarabinerRules() Task {
	return Task{
		name: "Unselect other Karabiner profiles",
		exec: func(i Installation) error {
			i.applyRules("main.json")
			i.applyRules("finder.json")
			return nil
		},
	}
}

func ApplyAppLauncherRules() Task {
	return Task{
		name: "Unselect other Karabiner profiles",
		exec: func(i Installation) error {
			switch strings.ToLower(i.AppLauncher) {
			case "spotlight":
				i.applyRules("spotlight.json")
			case "launchpad":
				i.applyRules("launchpad.json")
			case "alfred":
				{
					if i.Exists("Alfred 4.app") || i.Exists("Alfred 5.app") {

						i.applyRules("alfred.json")

						dirs, err := findMatchingDirs(i.ApplicationSupportDir()+"/Alfred/Alfred.alfredpreferences/preferences/local", "", "hotkey", "prefs.plist")

						if err != nil {
							return err
						}

						for _, e := range dirs {
							copyFileFromEmbedFS("alfred/prefs.plist", e)
						}
					} else {
						printColored(Yellow, fmt.Sprintf("Alfred app not found. Skipping..."))
					}
				}
			}

			return nil
		},
	}
}

func ApplyKeyboardLayoutRules() Task {
	return Task{
		name: "Unselect other Karabiner profiles",
		exec: func(i Installation) error {
			switch strings.ToLower(i.KeyboardLayout) {
			case "mac":
				jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >tmp && mv tmp %s", i.profileName, i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
				i.Run(jq)
			}
			return nil
		},
	}
}

func ApplyTerminalRules() Task {
	return Task{
		name: "Unselect other Karabiner profiles",
		exec: func(i Installation) error {
			switch strings.ToLower(i.Terminal) {
			case "default":
				i.applyRules("apple-terminal.json")
			case "iterm":
				if i.Exists("iTerm.app") {
					i.applyRules("iterm.json")
				} else {
					printColored(Yellow, fmt.Sprintf("iTerm app not found. Skipping..."))
				}
			case "warp":
				{
					if i.Exists("Warp.app") {
						i.applyRules("warp.json")
					} else {
						printColored(Yellow, fmt.Sprintf("Warp app not found. Skipping..."))
					}
				}
			}

			return nil
		},
	}
}

func ReformatKarabinerConfigFile() Task {
	return Task{
		name: "Reformat Karabiner config file",
		exec: func(i Installation) error {
			i.Run(fmt.Sprintf("jq '.' %s > tmp && mv tmp %s", i.KarabinerConfigFile(), i.KarabinerConfigFile()))
			return nil
		},
	}
}

func OpenKarabiner() Task {
	return Task{
		name: "Open Karabiner-Elements.app",
		exec: func(i Installation) error {
			i.Run("open -a Karabiner-Elements")
			return nil
		},
	}
}

func CopyIdeKeymaps() Task {
	return Task{
		name: "Install IDE keymaps",
		exec: func(i Installation) error {
			for _, ide := range i.Ides {
				name, _ := IdeKeymapByFullName(ide)
				i.installIdeKeymap(name)
			}
			return nil
		},
	}
}

func CloseRectangle() Task {
	return Task{
		name: "Close rectangle",
		exec: func(i Installation) error {
			i.Run("killall Rectangle")
			return nil
		},
	}
}

func CopyRectanglePreferences() Task {
	return Task{
		name: "Install Rectangle preferences",
		exec: func(i Installation) error {
			rectanglePlist := filepath.Join(i.PreferencesDir(), "com.knollsoft.Rectangle.plist")
			copyFileFromEmbedFS("rectangle/Settings.xml", rectanglePlist)

			plutilCmdRectangle := fmt.Sprintf("plutil -convert binary1 %s", rectanglePlist)
			i.Run(plutilCmdRectangle)
			i.Run("defaults read com.knollsoft.Rectangle.plist")
			return nil
		},
	}
}

func OpenRectangle() Task {
	return Task{
		name: "Open Rectangle.app",
		exec: func(i Installation) error {
			i.Run("open -a Rectangle")
			return nil
		},
	}
}

func CloseAltTab() Task {
	return Task{
		name: "Close AtlTab.app",
		exec: func(i Installation) error {
			i.Run("killall AltTab")
			return nil
		},
	}
}

func InstallAltTabPreferences() Task {
	return Task{
		name: "Install AltTab preferences",
		exec: func(i Installation) error {
			altTabPlist := filepath.Join(i.PreferencesDir(), "/com.lwouis.alt-tab-macos.plist")
			copyFileFromEmbedFS("alt-tab/Settings.xml", altTabPlist)

			// set up blacklist
			var mappedStrings []string
			for _, s := range i.Blacklist {
				mappedStrings = append(mappedStrings, fmt.Sprintf(`{"ignore":"0","bundleIdentifier":"%s","hide":"1"}`, s))
			}

			result := "[" + strings.Join(mappedStrings, ",") + "]"

			replaceWordInFile(altTabPlist, "_BLACKLIST_", result)

			plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", altTabPlist)
			i.Run(plutilCmd)

			i.Run("defaults read com.lwouis.alt-tab-macos.plist")
			return nil
		},
	}
}

func OpenAltTab() Task {
	return Task{
		name: "Open AtlTab.app",
		exec: func(i Installation) error {
			i.Run("open -a AltTab")
			return nil
		},
	}
}

func ApplySystemSettings() Task {
	return Task{
		name: "Apply system settings",
		exec: func(i Installation) error {
			optionsMap := make(map[string]bool)
			for _, value := range i.SystemSettings {
				optionsMap[strings.ToLower(value)] = true
			}

			if optionsMap["enable dock auto-hide (2s delay)"] {
				i.Run("defaults write com.apple.dock autohide -bool true")
				i.Run("defaults write com.apple.dock autohide-delay -float 2 && killall Dock")
			}

			if optionsMap[`change dock minimize animation to "scale"`] {
				i.Run(`defaults write com.apple.dock "mineffect" -string "scale" && killall Dock`)
			}

			if optionsMap["enable home & end keys"] {
				copyFileFromEmbedFS("system/DefaultKeyBinding.dict", filepath.Join(i.LibraryDir(), "/KeyBindings/DefaultKeyBinding.dict"))
			}

			if optionsMap["show hidden files in finder"] {
				i.Run("defaults write com.apple.finder AppleShowAllFiles -bool true")
			}

			if optionsMap["show directories on top in finder"] {
				i.Run("defaults write com.apple.finder _FXSortFoldersFirst -bool true")
			}
			if optionsMap["show full posix paths in finder window title"] {
				i.Run("defaults write com.apple.finder _FXShowPosixPathInTitle -bool true")
			}
			return nil
		},
	}
}

func (i Installation) applyRules(file string) {
	copyFileFromEmbedFS(filepath.Join("karabiner", file), filepath.Join(i.KarabinerComplexModificationsDir(), file))
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/%s >tmp && mv tmp %s", i.profileName, i.profileName, i.KarabinerConfigFile(), i.KarabinerComplexModificationsDir(), file, i.KarabinerConfigFile())
	i.Run(jq)
}

func (i Installation) installIdeKeymap(ide IDE) error {

	var destDirs []string

	if ide.multipleDirs {
		destDirs = i.IdeKeymapPaths(ide)
	} else {
		destDirs = []string{filepath.Join(i.Path, ide.parentDir, ide.dir, ide.keymapsDir, ide.destKeymapsFile)}
	}

	if len(destDirs) == 0 {
		printColored(Yellow, fmt.Sprintf("%s not found. Skipping...", ide.fullName))
		return nil
	}

	for _, d := range destDirs {
		err := copyFileFromEmbedFS(i.SourceKeymap(ide), d)

		if err != nil {
			return err
		}
	}

	return nil
}
