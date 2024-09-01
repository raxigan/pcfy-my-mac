package task

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"os"
	"path/filepath"
	"strings"
)

type Task struct {
	Name    string
	Execute func(i install.Installation) error
}

func DownloadDependencies() Task {
	return Task{
		Name: "Install dependencies",
		Execute: func(i install.Installation) error {

			var notInstalled []string
			var commands []string

			all := []Dependency{JqDependency(), KarabinerDependency(), AltTabDependency(), RectangleDependency()}

			for _, d := range all {
				if !common.Exists(d.command) {
					notInstalled = append(notInstalled, d.name)
					commands = append(commands, d.installCommand)
				}
			}

			if len(notInstalled) > 0 {
				installApp := false
				prompt := &survey.Confirm{
					Message: fmt.Sprintf("The following dependencies will be installed: %s. Do you agree?", strings.Join(notInstalled, ", ")),
				}
				common.HandleInterrupt(
					survey.AskOne(prompt, &installApp),
				)

				if !installApp {
					fmt.Printf("Qutting...")
					i.Exit(0)
				}

				for _, c := range commands {
					i.Run(c)
				}
			}

			return nil

		},
	}
}

func CloseKarabiner() Task {
	return Task{
		Name:    "Close Karabiner",
		Execute: func(i install.Installation) error { i.Run("killall Karabiner-Menu"); return nil },
	}
}

func BackupKarabinerConfig() Task {
	return Task{
		Name: "Backup karabiner config",
		Execute: func(i install.Installation) error {
			original := i.KarabinerConfigFile()
			backupDest := i.KarabinerConfigBackupFile(i.InstallationTime)

			configExists := common.FileExists(original)

			if !configExists {
				os.MkdirAll(i.KarabinerComplexModificationsDir(), 0755)
				copyFile(filepath.Join("karabiner", "default.json"), original, i)
			}

			err := common.CopyFile(original, backupDest)

			if err != nil {
				return err
			}

			return nil
		},
	}
}

func DeleteExistingKarabinerProfile() Task {
	return Task{
		Name: "Delete existing Karabiner profile",
		Execute: func(i install.Installation) error {
			deleteProfileJqCmd := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >tmp && mv tmp %s", i.ProfileName, i.ProfileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(deleteProfileJqCmd)
			return nil
		},
	}
}

func CreateKarabinerProfile() Task {
	return Task{
		Name: "Create new Karabiner profile",
		Execute: func(i install.Installation) error {
			copyFile("karabiner/karabiner-profile.json", "tmp", i)
			addProfileJqCmd := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile tmp --indent 4 >INPUT.tmp && mv INPUT.tmp %s && rm tmp", i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(addProfileJqCmd)
			return nil
		},
	}
}

func NameKarabinerProfile() Task {
	return Task{
		Name: "Rename new Karabiner profile",
		Execute: func(i install.Installation) error {
			copyFile("karabiner/karabiner-profile.json", "tmp", i)
			addProfileJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name == \"_PROFILE_NAME_\" then .name = \"%s\" else . end)' %s > tmp && mv tmp %s", i.ProfileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(addProfileJqCmd)
			return nil
		},
	}
}

func UnselectOtherKarabinerProfiles() Task {
	return Task{
		Name: "Unselect other Karabiner profiles",
		Execute: func(i install.Installation) error {
			unselectJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name != \"%s\" then .selected = false else . end)' %s > tmp && mv tmp %s", i.ProfileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
			i.Run(unselectJqCmd)
			return nil
		},
	}
}

func ApplyMainKarabinerRules() Task {
	return Task{
		Name: "Apply main Karabiner rules",
		Execute: func(i install.Installation) error {
			ApplyRules(i, "main.json")
			ApplyRules(i, "finder.json")
			return nil
		},
	}
}

func ApplyAppLauncherRules() Task {
	return Task{
		Name: "Apply app launcher rules",
		Execute: func(i install.Installation) error {
			switch strings.ToLower(i.AppLauncher) {
			case strings.ToLower(param.None):
				ApplyRules(i, "app-launcher-none.json")
			case strings.ToLower(param.Spotlight):
				ApplyRules(i, "spotlight.json")
			case strings.ToLower(param.Launchpad):
				ApplyRules(i, "launchpad.json")
			case strings.ToLower(param.Alfred):
				{
					if common.Exists("Alfred 4.app") || common.Exists("Alfred 5.app") {

						ApplyRules(i, "alfred.json")

						paths, err := common.FindMatchingPaths(i.ApplicationSupportDir()+"/Alfred/Alfred.alfredpreferences/preferences/local/{version}/hotkey", "prefs.plist")

						if err != nil {
							return err
						}

						for _, path := range paths {
							copyFile("alfred/prefs.plist", path, i)
						}
					} else {
						i.TryLog(install.WarnMsg, "Alfred app not found. Skipping...")
					}
				}
			default:
				return errors.New("Unknown app launcher: " + i.AppLauncher)
			}

			return nil
		},
	}
}

func ApplyKeyboardLayoutRules() Task {
	return Task{
		Name: "Apply keyboard layout rules",
		Execute: func(i install.Installation) error {
			switch strings.ToLower(i.KeyboardLayout) {
			case strings.ToLower(param.PC):
			case strings.ToLower(param.Mac), strings.ToLower(param.None):
				jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >tmp && mv tmp %s", i.ProfileName, i.ProfileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
				i.Run(jq)
			default:
				return errors.New("Unknown keyboard layout: " + i.KeyboardLayout)
			}
			return nil
		},
	}
}

func ApplyTerminalRules() Task {
	return Task{
		Name: "Apply terminal rules",
		Execute: func(i install.Installation) error {

			switch strings.ToLower(i.Terminal) {
			case strings.ToLower(param.None):
			case strings.ToLower(param.Default):
				ApplyRules(i, "apple-terminal.json")
			case strings.ToLower(param.ITerm):
				if common.Exists("iTerm.app") {
					ApplyRules(i, "iterm.json")
				} else {
					common.PrintColored(common.Yellow, fmt.Sprintf("iTerm app not found. Skipping..."))
				}
			case strings.ToLower(param.Warp):
				{
					if common.Exists("Warp.app") {
						ApplyRules(i, "warp.json")
					} else {
						common.PrintColored(common.Yellow, fmt.Sprintf("Warp app not found. Skipping..."))
					}
				}
			case strings.ToLower(param.Wave):
				{
					if common.Exists("Wave.app") {
						ApplyRules(i, "wave.json")
					} else {
						common.PrintColored(common.Yellow, fmt.Sprintf("Wave app not found. Skipping..."))
					}
				}
			default:
				return errors.New("Unknown terminal: " + i.Terminal)
			}

			return nil
		},
	}
}

func ReformatKarabinerConfigFile() Task {
	return Task{
		Name: "Reformat Karabiner config file",
		Execute: func(i install.Installation) error {
			i.Run(fmt.Sprintf("jq '.' %s > tmp && mv tmp %s", i.KarabinerConfigFile(), i.KarabinerConfigFile()))
			return nil
		},
	}
}

func OpenKarabiner() Task {
	return Task{
		Name: "Open Karabiner-Elements.app",
		Execute: func(i install.Installation) error {
			i.Run("open -a Karabiner-Elements")
			return nil
		},
	}
}

func CopyIdeKeymaps() Task {
	return Task{
		Name: "Install IDE keymaps",
		Execute: func(i install.Installation) error {
			for _, keymap := range i.Keymaps {
				name, _ := param.IdeKeymapByFullName(keymap)
				InstallIdeKeymap(i, name)
			}
			return nil
		},
	}
}

func CloseRectangle() Task {
	return Task{
		Name: "Close rectangle",
		Execute: func(i install.Installation) error {
			i.Run("killall Rectangle")
			return nil
		},
	}
}

func CopyRectanglePreferences() Task {
	return Task{
		Name: "Install Rectangle preferences",
		Execute: func(i install.Installation) error {
			rectanglePlist := filepath.Join(i.PreferencesDir(), "com.knollsoft.Rectangle.plist")
			copyFile("rectangle/com.knollsoft.Rectangle.plist", rectanglePlist, i)

			plutilCmdRectangle := fmt.Sprintf("plutil -convert binary1 %s", rectanglePlist)
			i.Run(plutilCmdRectangle)
			i.Run("defaults read com.knollsoft.Rectangle.plist")
			return nil
		},
	}
}

func OpenRectangle() Task {
	return Task{
		Name: "Open Rectangle.app",
		Execute: func(i install.Installation) error {
			i.Run("open -a Rectangle")
			return nil
		},
	}
}

func CloseAltTab() Task {
	return Task{
		Name: "Close AtlTab.app",
		Execute: func(i install.Installation) error {
			i.Run("killall AltTab")
			return nil
		},
	}
}

func InstallAltTabPreferences() Task {
	return Task{
		Name: "Install AltTab preferences",
		Execute: func(i install.Installation) error {

			i.Commander.TryLog(install.TaskMsg, fmt.Sprintf("Exclude %s from AltTab", i.Blacklist))

			altTabPlist := filepath.Join(i.PreferencesDir(), "com.lwouis.alt-tab-macos.plist")
			copyFile("alt-tab/com.lwouis.alt-tab-macos.plist", altTabPlist, i)

			var mappedStrings []string
			for _, bundle := range i.Blacklist {
				mappedStrings = append(mappedStrings, fmt.Sprintf(`{"ignore":"0","bundleIdentifier":"%s","hide":"1"}`, bundle))
			}

			result := "[" + strings.Join(mappedStrings, ",") + "]"

			common.ReplaceWordInFile(altTabPlist, "_BLACKLIST_", result)

			plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", altTabPlist)
			i.Run(plutilCmd)

			i.Run("defaults read com.lwouis.alt-tab-macos.plist")
			return nil
		},
	}
}

func OpenAltTab() Task {
	return Task{
		Name: "Open AtlTab.app",
		Execute: func(i install.Installation) error {
			i.Run("open -a AltTab")
			return nil
		},
	}
}

func ApplySystemSettings() Task {
	return Task{
		Name: "Apply system settings",
		Execute: func(i install.Installation) error {
			for _, value := range i.SystemSettings {
				simpleParamName := param.ToSimpleParamName(value)

				switch simpleParamName {
				case "enable-dock-auto-hide-2s-delay":
					{
						i.Run("defaults write com.apple.dock autohide -bool true")
						i.Run("defaults write com.apple.dock autohide-delay -float 2 && killall Dock")
					}
				case "change-dock-minimize-animation-to-scale":
					{
						i.Run(`defaults write com.apple.dock "mineffect" -string "scale" && killall Dock`)
					}
				case "enable-home-and-end-keys":
					{
						copyFile("system/DefaultKeyBinding.dict", filepath.Join(i.LibraryDir(), "/KeyBindings/DefaultKeyBinding.dict"), i)
					}
				case "show-hidden-files-in-finder":
					{
						i.Run("defaults write com.apple.finder AppleShowAllFiles -bool true")
					}
				case "show-directories-on-top-in-finder":
					{
						i.Run("defaults write com.apple.finder _FXSortFoldersFirst -bool true")
					}
				case "show-full-posix-paths-in-finder-window-title":
					{
						i.Run("defaults write com.apple.finder _FXShowPosixPathInTitle -bool true")
					}

				}
			}

			return nil
		},
	}
}

func ApplyRules(i install.Installation, file string) {
	copyFile(filepath.Join("karabiner", file), filepath.Join(i.KarabinerComplexModificationsDir(), file), i)
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/%s >tmp && mv tmp %s", i.ProfileName, i.ProfileName, i.KarabinerConfigFile(), i.KarabinerComplexModificationsDir(), file, i.KarabinerConfigFile())
	i.Run(jq)
}

func InstallIdeKeymap(i install.Installation, ide param.IDE) error {

	var destDirs = i.IdeKeymapPaths(ide)

	if len(destDirs) == 0 {
		i.TryLog(install.WarnMsg, fmt.Sprintf("%s not found. Skipping...", ide.FullName))
		return nil
	}

	for _, d := range destDirs {
		err := copyFile(i.SourceKeymap(ide), d, i)

		if err != nil {
			return err
		}
	}

	return nil
}

func CopyHidutilRemappingFile() Task {
	return Task{
		Name: "Copy hidutil remapping file",
		Execute: func(i install.Installation) error {
			copyFile("system/com.github.pcfy-my-mac.plist", filepath.Join(i.LaunchAgents(), "com.github.pcfy-my-mac.plist"), i)
			return nil
		},
	}
}

func ExecuteHidutil() Task {
	return Task{
		Name: "Execute hidutil command",
		Execute: func(i install.Installation) error {
			remappingFile, _ := common.ReadFileFromEmbedFS("system/com.github.pcfy-my-mac.plist")
			start := strings.Index(remappingFile, "<array>")
			end := strings.Index(remappingFile, "</array>")
			arrayContent := remappingFile[start+len("<array>") : end]

			arrayContent = strings.ReplaceAll(arrayContent, "<string>{\"", "<string>'{\"")
			arrayContent = strings.ReplaceAll(arrayContent, "]}</string>", "]}'</string>")

			arrayContent = strings.ReplaceAll(arrayContent, "<string>", "")
			arrayContent = strings.ReplaceAll(arrayContent, "</string>", "")

			command := strings.TrimSpace(arrayContent)
			command = strings.Join(strings.Fields(command), " ")
			command = strings.ReplaceAll(command, "/usr/bin/hidutil", "hidutil")

			i.Run(command)

			return nil
		},
	}
}

func copyFile(src, dst string, i install.Installation) error {
	loggedSrc := strings.ReplaceAll(src, i.HomeDir.Path, "~")
	loggedDst := strings.ReplaceAll(dst, i.HomeDir.Path, "~")
	i.TryLog(install.FileMsg, fmt.Sprintf("Copy file %s to %s", loggedSrc, loggedDst))
	return common.CopyFileFromEmbedFS(src, dst)
}
